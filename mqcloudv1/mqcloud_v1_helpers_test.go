//go:build integration || examples

/**
 * (C) Copyright IBM Corp. 2026.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package mqcloudv1_test

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/IBM/mqcloud-go-sdk/mqcloudv1"
	. "github.com/onsi/ginkgo"
)

// Helpers shared by the integration and examples tests for waiting on asynchronous APIs.
// This file is hand maintained, it is not produced by the generator.

const (
	asyncPollInterval = 5 * time.Second
	asyncTimeout      = 10 * time.Minute
)

// WaitForOperation polls GetOperation until the operation has succeeded.
// It returns an error if the operation fails, if it cannot be read, or if it does not finish before the timeout.
func WaitForOperation(service *mqcloudv1.MqcloudV1, serviceInstanceGuid string, operationID string) error {
	deadline := time.Now().Add(asyncTimeout)
	for {
		operationStatus, _, err := service.GetOperation(&mqcloudv1.GetOperationOptions{
			ServiceInstanceGuid: &serviceInstanceGuid,
			OperationID:         &operationID,
		})
		if err != nil {
			return fmt.Errorf("unable to get operation %s: %w", operationID, err)
		}

		status := ""
		if operationStatus.Status != nil {
			status = *operationStatus.Status
		}
		switch status {
		case mqcloudv1.OperationStatus_Status_Succeeded:
			return nil
		case mqcloudv1.OperationStatus_Status_Failed:
			var messages []string
			for _, operationError := range operationStatus.Errors {
				code, message := "", ""
				if operationError.Code != nil {
					code = *operationError.Code
				}
				if operationError.Message != nil {
					message = *operationError.Message
				}
				messages = append(messages, code+": "+message)
			}
			return fmt.Errorf("operation %s failed: %s", operationID, strings.Join(messages, "; "))
		}

		if time.Now().After(deadline) {
			return fmt.Errorf("operation %s did not complete within %s, last status %q", operationID, asyncTimeout, status)
		}
		time.Sleep(asyncPollInterval)
	}
}

// WaitForQueueManagerRunning polls GetQueueManagerStatus until the queue manager is running.
// It returns an error if the queue manager reaches a failed state, or does not start before the timeout.
func WaitForQueueManagerRunning(service *mqcloudv1.MqcloudV1, serviceInstanceGuid string, queueManagerID string) error {
	deadline := time.Now().Add(asyncTimeout)
	for {
		queueManagerStatus, _, err := service.GetQueueManagerStatus(&mqcloudv1.GetQueueManagerStatusOptions{
			ServiceInstanceGuid: &serviceInstanceGuid,
			QueueManagerID:      &queueManagerID,
		})
		// The queue manager may not be visible immediately after the create request, so keep polling on error.
		status := ""
		if err == nil && queueManagerStatus.Status != nil {
			status = *queueManagerStatus.Status
		}
		switch status {
		case mqcloudv1.QueueManagerStatus_Status_Running:
			return nil
		case mqcloudv1.QueueManagerStatus_Status_Failed,
			mqcloudv1.QueueManagerStatus_Status_InitializationFailed,
			mqcloudv1.QueueManagerStatus_Status_RestoreFailed:
			return fmt.Errorf("queue manager %s reached state %q", queueManagerID, status)
		}

		if time.Now().After(deadline) {
			return fmt.Errorf("queue manager %s was not running within %s, last status %q, last error %v", queueManagerID, asyncTimeout, status, err)
		}
		time.Sleep(asyncPollInterval)
	}
}

// SkipUnlessConfigured skips the current spec unless every named config value is set. It is the one place the tests
// decide whether to skip. A value is set either in the env file (as MQCLOUD_<NAME>) or by an earlier test that captures
// it, so this also skips the tests that depend on a resource when the test that creates it did not run or failed.
// For example VERSION_UPGRADE is only set by the GetOptions test when a newer queue manager version is offered.
func SkipUnlessConfigured(config map[string]string, keys ...string) {
	var missing []string
	for _, key := range keys {
		if config[key] == "" {
			missing = append(missing, "MQCLOUD_"+key)
		}
	}
	if len(missing) > 0 {
		Skip("Skipping, not set in mqcloud_v1.env or by an earlier test: " + strings.Join(missing, ", "))
	}
}

// SkipTestIfQmIsNotRunning skips the current spec unless the queue manager exists and is running.
// It also skips unless the deployment ID and queue manager ID are set, so callers do not need to check those first.
// An error from the status call fails the spec, so an API problem is not hidden by a skip.
func SkipTestIfQmIsNotRunning(config map[string]string, service *mqcloudv1.MqcloudV1) {
	SkipUnlessConfigured(config, "SERVICE_INSTANCE_DEPLOYMENT_GUID", "QUEUE_MANAGER_ID")
	if service == nil {
		Skip("Skipping, the service client has not been created")
	}
	deploymentID, qmID := config["SERVICE_INSTANCE_DEPLOYMENT_GUID"], config["QUEUE_MANAGER_ID"]
	queueManagerStatus, _, err := service.GetQueueManagerStatus(&mqcloudv1.GetQueueManagerStatusOptions{
		ServiceInstanceGuid: &deploymentID,
		QueueManagerID:      &qmID,
	})
	if err != nil {
		Fail(fmt.Sprintf("unable to get the status of queue manager %s: %s", qmID, err.Error()))
	}
	if queueManagerStatus.Status == nil || *queueManagerStatus.Status != mqcloudv1.QueueManagerStatus_Status_Running {
		status := "unknown"
		if queueManagerStatus.Status != nil {
			status = *queueManagerStatus.Status
		}
		Skip(fmt.Sprintf("Skipping, queue manager %s is %s, not running", qmID, status))
	}
}

// CompareVersions compares two queue manager versions such as "9.4.0_1" and "10.0.0_3".
// It returns a negative number if a is older than b, zero if they are equal and a positive number if a is newer.
// The numbers are compared numerically, so "10.0.0_1" is newer than "9.4.0_5".
func CompareVersions(a string, b string) int {
	partsA, partsB := versionParts(a), versionParts(b)
	for i := 0; i < len(partsA) || i < len(partsB); i++ {
		x, y := 0, 0
		if i < len(partsA) {
			x = partsA[i]
		}
		if i < len(partsB) {
			y = partsB[i]
		}
		if x != y {
			return x - y
		}
	}
	return strings.Compare(a, b)
}

func versionParts(version string) []int {
	var parts []int
	for _, field := range strings.FieldsFunc(version, func(r rune) bool { return r == '.' || r == '_' }) {
		number, err := strconv.Atoi(field)
		if err != nil {
			number = 0
		}
		parts = append(parts, number)
	}
	return parts
}

// SortVersions returns a copy of the versions ordered from oldest to newest.
func SortVersions(versions []string) []string {
	sorted := append([]string(nil), versions...)
	sort.SliceStable(sorted, func(i, j int) bool { return CompareVersions(sorted[i], sorted[j]) < 0 })
	return sorted
}

// WaitForQueueManagerVersion polls GetQueueManager until the queue manager reports the given version,
// then waits for it to be running again. Use it after SetQueueManagerVersion, where the queue manager can
// briefly still report running at the old version.
func WaitForQueueManagerVersion(service *mqcloudv1.MqcloudV1, serviceInstanceGuid string, queueManagerID string, version string) error {
	deadline := time.Now().Add(asyncTimeout)
	for {
		queueManagerDetails, _, err := service.GetQueueManager(&mqcloudv1.GetQueueManagerOptions{
			ServiceInstanceGuid: &serviceInstanceGuid,
			QueueManagerID:      &queueManagerID,
		})
		current := ""
		if err == nil && queueManagerDetails.Version != nil {
			current = *queueManagerDetails.Version
		}
		if current == version {
			return WaitForQueueManagerRunning(service, serviceInstanceGuid, queueManagerID)
		}

		if time.Now().After(deadline) {
			return fmt.Errorf("queue manager %s did not reach version %q within %s, last version %q, last error %v", queueManagerID, version, asyncTimeout, current, err)
		}
		time.Sleep(asyncPollInterval)
	}
}

// ChooseQueueManagerVersions picks the version to create a queue manager at, and the version to upgrade it to.
// A value that is already set (MQCLOUD_VERSION or MQCLOUD_VERSION_UPGRADE in the env file) is kept as it is, so someone
// can control the versions used. Whatever is not set is chosen from the versions on offer: the queue manager is created
// at the oldest and upgraded to the newest. The upgrade version stays empty when nothing newer than the creation version
// is on offer, which makes the upgrade test skip.
func ChooseQueueManagerVersions(offered []string, version string, upgradeVersion string) (string, string) {
	sorted := SortVersions(offered)
	if version == "" && len(sorted) > 0 {
		version = sorted[0]
	}
	if upgradeVersion == "" && len(sorted) > 0 && CompareVersions(sorted[len(sorted)-1], version) > 0 {
		upgradeVersion = sorted[len(sorted)-1]
	}
	return version, upgradeVersion
}
