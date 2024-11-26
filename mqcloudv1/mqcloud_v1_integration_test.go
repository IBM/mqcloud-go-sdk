//go:build integration

/**
 * (C) Copyright IBM Corp. 2024.
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
	"log"
	"os"
	"time"

	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/mqcloud-go-sdk/mqcloudv1"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

/**
 * This file contains an integration test for the mqcloudv1 package.
 *
 * Notes:
 *
 * The integration test will automatically skip tests if the required config file is not available.
 */

var _ = Describe(`MqcloudV1 Integration Tests`, func() {
	const externalConfigFile = "../mqcloud_v1.env"

	var (
		err            error
		mqcloudService *mqcloudv1.MqcloudV1
		serviceURL     string
		config         map[string]string
	)

	var shouldSkipTest = func() {
		Skip("External configuration is not available, skipping tests...")
	}

	Describe(`External configuration`, func() {
		It("Successfully load the configuration", func() {
			_, err = os.Stat(externalConfigFile)
			if err != nil {
				Skip("External configuration file not found, skipping tests: " + err.Error())
			}

			os.Setenv("IBM_CREDENTIALS_FILE", externalConfigFile)
			config, err = core.GetServiceProperties(mqcloudv1.DefaultServiceName)
			if err != nil {
				Skip("Error loading service properties, skipping tests: " + err.Error())
			}
			serviceURL = config["URL"]
			if serviceURL == "" {
				Skip("Unable to load service URL configuration property, skipping tests")
			}

			fmt.Fprintf(GinkgoWriter, "Service URL: %v\n", serviceURL)
			shouldSkipTest = func() {}
		})
	})

	Describe(`Client initialization`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It("Successfully construct the service client instance", func() {
			mqcloudServiceOptions := &mqcloudv1.MqcloudV1Options{}

			mqcloudService, err = mqcloudv1.NewMqcloudV1UsingExternalConfig(mqcloudServiceOptions)
			Expect(err).To(BeNil())
			Expect(mqcloudService).ToNot(BeNil())
			Expect(mqcloudService.Service.Options.URL).To(Equal(serviceURL))

			core.SetLogger(core.NewLogger(core.LevelDebug, log.New(GinkgoWriter, "", log.LstdFlags), log.New(GinkgoWriter, "", log.LstdFlags)))
			mqcloudService.EnableRetries(4, 30*time.Second)
		})
	})

	Describe(`GetUsageDetails - Get the usage details`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`GetUsageDetails(getUsageDetailsOptions *GetUsageDetailsOptions)`, func() {
			getUsageDetailsOptions := &mqcloudv1.GetUsageDetailsOptions{
				ServiceInstanceGuid: core.StringPtr(config["SERVICE_INSTANCE_DEPLOYMENT_GUID"]),
			}

			usage, response, err := mqcloudService.GetUsageDetails(getUsageDetailsOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(usage).ToNot(BeNil())
		})
	})

	Describe(`GetOptions - Return configuration options (eg, available deployment locations, queue manager sizes)`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`GetOptions(getOptionsOptions *GetOptionsOptions)`, func() {
			getOptionsOptions := &mqcloudv1.GetOptionsOptions{
				ServiceInstanceGuid: core.StringPtr(config["SERVICE_INSTANCE_DEPLOYMENT_GUID"]),
			}

			configurationOptions, response, err := mqcloudService.GetOptions(getOptionsOptions)
			Expect(err).To(BeNil())
			config["LOCATION"] = configurationOptions.Locations[0]
			Expect(response.StatusCode).To(Equal(200))
			Expect(configurationOptions).ToNot(BeNil())
		})
	})

	Describe(`CreateQueueManager - Create a new queue manager`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`CreateQueueManager(createQueueManagerOptions *CreateQueueManagerOptions)`, func() {
			createQueueManagerOptions := &mqcloudv1.CreateQueueManagerOptions{
				ServiceInstanceGuid: core.StringPtr(config["SERVICE_INSTANCE_DEPLOYMENT_GUID"]),
				Name:                core.StringPtr("testqm10"),
				Location:            core.StringPtr(config["LOCATION"]),
				Size:                core.StringPtr("xsmall"),
				DisplayName:         core.StringPtr("A test queue manager"),
				Version:             core.StringPtr(config["VERSION"]),
			}

			queueManagerTaskStatus, response, err := mqcloudService.CreateQueueManager(createQueueManagerOptions)
			Expect(err).To(BeNil())
			config["QUEUE_MANAGER_ID"] = *queueManagerTaskStatus.QueueManagerID
			Expect(response.StatusCode).To(Equal(202))
			Expect(queueManagerTaskStatus).ToNot(BeNil())
		})
	})

	Describe(`ListQueueManagers - Get list of queue managers`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`ListQueueManagers(listQueueManagersOptions *ListQueueManagersOptions) with pagination`, func() {
			listQueueManagersOptions := &mqcloudv1.ListQueueManagersOptions{
				ServiceInstanceGuid: core.StringPtr(config["SERVICE_INSTANCE_DEPLOYMENT_GUID"]),
				Offset:              core.Int64Ptr(int64(0)),
				Limit:               core.Int64Ptr(int64(10)),
			}

			listQueueManagersOptions.Offset = nil
			listQueueManagersOptions.Limit = core.Int64Ptr(1)

			var allResults []mqcloudv1.QueueManagerDetails
			for {
				queueManagerDetailsCollection, response, err := mqcloudService.ListQueueManagers(listQueueManagersOptions)
				Expect(err).To(BeNil())
				Expect(response.StatusCode).To(Equal(200))
				Expect(queueManagerDetailsCollection).ToNot(BeNil())
				allResults = append(allResults, queueManagerDetailsCollection.QueueManagers...)

				listQueueManagersOptions.Offset, err = queueManagerDetailsCollection.GetNextOffset()
				Expect(err).To(BeNil())

				if listQueueManagersOptions.Offset == nil {
					break
				}
			}
			fmt.Fprintf(GinkgoWriter, "Retrieved a total of %d item(s) with pagination.\n", len(allResults))
		})
		It(`ListQueueManagers(listQueueManagersOptions *ListQueueManagersOptions) using QueueManagersPager`, func() {
			listQueueManagersOptions := &mqcloudv1.ListQueueManagersOptions{
				ServiceInstanceGuid: core.StringPtr(config["SERVICE_INSTANCE_DEPLOYMENT_GUID"]),
				Limit:               core.Int64Ptr(int64(10)),
			}

			// Test GetNext().
			pager, err := mqcloudService.NewQueueManagersPager(listQueueManagersOptions)
			Expect(err).To(BeNil())
			Expect(pager).ToNot(BeNil())

			var allResults []mqcloudv1.QueueManagerDetails
			for pager.HasNext() {
				nextPage, err := pager.GetNext()
				Expect(err).To(BeNil())
				Expect(nextPage).ToNot(BeNil())
				allResults = append(allResults, nextPage...)
			}

			// Test GetAll().
			pager, err = mqcloudService.NewQueueManagersPager(listQueueManagersOptions)
			Expect(err).To(BeNil())
			Expect(pager).ToNot(BeNil())

			allItems, err := pager.GetAll()
			Expect(err).To(BeNil())
			Expect(allItems).ToNot(BeNil())

			Expect(len(allItems)).To(Equal(len(allResults)))
			fmt.Fprintf(GinkgoWriter, "ListQueueManagers() returned a total of %d item(s) using QueueManagersPager.\n", len(allResults))
		})
	})

	Describe(`GetQueueManager - Get details of a queue manager`, func() {
		BeforeEach(func() {
			shouldSkipTest()
			SkipTestIfQmIsNotRunning(config["QUEUE_MANAGER_ID"], mqcloudService, config["SERVICE_INSTANCE_DEPLOYMENT_GUID"])
		})
		It(`GetQueueManager(getQueueManagerOptions *GetQueueManagerOptions)`, func() {
			getQueueManagerOptions := &mqcloudv1.GetQueueManagerOptions{
				ServiceInstanceGuid: core.StringPtr(config["SERVICE_INSTANCE_DEPLOYMENT_GUID"]),
				QueueManagerID:      core.StringPtr(config["QUEUE_MANAGER_ID"]),
			}

			queueManagerDetails, response, err := mqcloudService.GetQueueManager(getQueueManagerOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(queueManagerDetails).ToNot(BeNil())
		})
	})

	Describe(`SetQueueManagerVersion - Upgrade a queue manager`, func() {
		BeforeEach(func() {
			shouldSkipTest()
			SkipTestIfQmIsNotRunning(config["QUEUE_MANAGER_ID"], mqcloudService, config["SERVICE_INSTANCE_DEPLOYMENT_GUID"])
		})
		It(`SetQueueManagerVersion(setQueueManagerVersionOptions *SetQueueManagerVersionOptions)`, func() {
			setQueueManagerVersionOptions := &mqcloudv1.SetQueueManagerVersionOptions{
				ServiceInstanceGuid: core.StringPtr(config["SERVICE_INSTANCE_DEPLOYMENT_GUID"]),
				QueueManagerID:      core.StringPtr(config["QUEUE_MANAGER_ID"]),
				Version:             core.StringPtr(config["VERSION_UPGRADE"]),
			}

			queueManagerTaskStatus, response, err := mqcloudService.SetQueueManagerVersion(setQueueManagerVersionOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(202))
			Expect(queueManagerTaskStatus).ToNot(BeNil())
		})
	})

	Describe(`GetQueueManagerAvailableUpgradeVersions - Get the list of available versions that this queue manager can be upgraded to`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`GetQueueManagerAvailableUpgradeVersions(getQueueManagerAvailableUpgradeVersionsOptions *GetQueueManagerAvailableUpgradeVersionsOptions)`, func() {
			getQueueManagerAvailableUpgradeVersionsOptions := &mqcloudv1.GetQueueManagerAvailableUpgradeVersionsOptions{
				ServiceInstanceGuid: core.StringPtr(config["SERVICE_INSTANCE_DEPLOYMENT_GUID"]),
				QueueManagerID:      core.StringPtr(config["QUEUE_MANAGER_ID"]),
			}

			queueManagerVersionUpgrades, response, err := mqcloudService.GetQueueManagerAvailableUpgradeVersions(getQueueManagerAvailableUpgradeVersionsOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(queueManagerVersionUpgrades).ToNot(BeNil())
		})
	})

	Describe(`GetQueueManagerConnectionInfo - Get connection information for a queue manager`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`GetQueueManagerConnectionInfo(getQueueManagerConnectionInfoOptions *GetQueueManagerConnectionInfoOptions)`, func() {
			getQueueManagerConnectionInfoOptions := &mqcloudv1.GetQueueManagerConnectionInfoOptions{
				ServiceInstanceGuid: core.StringPtr(config["SERVICE_INSTANCE_DEPLOYMENT_GUID"]),
				QueueManagerID:      core.StringPtr(config["QUEUE_MANAGER_ID"]),
			}

			connectionInfo, response, err := mqcloudService.GetQueueManagerConnectionInfo(getQueueManagerConnectionInfoOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(connectionInfo).ToNot(BeNil())
		})
	})

	Describe(`GetQueueManagerStatus - Get the status of the queue manager`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`GetQueueManagerStatus(getQueueManagerStatusOptions *GetQueueManagerStatusOptions)`, func() {
			getQueueManagerStatusOptions := &mqcloudv1.GetQueueManagerStatusOptions{
				ServiceInstanceGuid: core.StringPtr(config["SERVICE_INSTANCE_DEPLOYMENT_GUID"]),
				QueueManagerID:      core.StringPtr(config["QUEUE_MANAGER_ID"]),
			}

			queueManagerStatus, response, err := mqcloudService.GetQueueManagerStatus(getQueueManagerStatusOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(queueManagerStatus).ToNot(BeNil())
		})
	})

	Describe(`ListUsers - Get a list of users for an instance`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`ListUsers(listUsersOptions *ListUsersOptions) with pagination`, func() {
			listUsersOptions := &mqcloudv1.ListUsersOptions{
				ServiceInstanceGuid: core.StringPtr(config["SERVICE_INSTANCE_DEPLOYMENT_GUID"]),
				Offset:              core.Int64Ptr(int64(0)),
				Limit:               core.Int64Ptr(int64(10)),
			}

			listUsersOptions.Offset = nil
			listUsersOptions.Limit = core.Int64Ptr(1)

			var allResults []mqcloudv1.UserDetails
			for {
				userDetailsCollection, response, err := mqcloudService.ListUsers(listUsersOptions)
				Expect(err).To(BeNil())
				Expect(response.StatusCode).To(Equal(200))
				Expect(userDetailsCollection).ToNot(BeNil())
				allResults = append(allResults, userDetailsCollection.Users...)

				listUsersOptions.Offset, err = userDetailsCollection.GetNextOffset()
				Expect(err).To(BeNil())

				if listUsersOptions.Offset == nil {
					break
				}
			}
			fmt.Fprintf(GinkgoWriter, "Retrieved a total of %d item(s) with pagination.\n", len(allResults))
		})
		It(`ListUsers(listUsersOptions *ListUsersOptions) using UsersPager`, func() {
			listUsersOptions := &mqcloudv1.ListUsersOptions{
				ServiceInstanceGuid: core.StringPtr(config["SERVICE_INSTANCE_DEPLOYMENT_GUID"]),
				Limit:               core.Int64Ptr(int64(10)),
			}

			// Test GetNext().
			pager, err := mqcloudService.NewUsersPager(listUsersOptions)
			Expect(err).To(BeNil())
			Expect(pager).ToNot(BeNil())

			var allResults []mqcloudv1.UserDetails
			for pager.HasNext() {
				nextPage, err := pager.GetNext()
				Expect(err).To(BeNil())
				Expect(nextPage).ToNot(BeNil())
				allResults = append(allResults, nextPage...)
			}

			// Test GetAll().
			pager, err = mqcloudService.NewUsersPager(listUsersOptions)
			Expect(err).To(BeNil())
			Expect(pager).ToNot(BeNil())

			allItems, err := pager.GetAll()
			Expect(err).To(BeNil())
			Expect(allItems).ToNot(BeNil())

			Expect(len(allItems)).To(Equal(len(allResults)))
			fmt.Fprintf(GinkgoWriter, "ListUsers() returned a total of %d item(s) using UsersPager.\n", len(allResults))
		})
	})

	Describe(`CreateUser - Add a user to an instance`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`CreateUser(createUserOptions *CreateUserOptions)`, func() {
			createUserOptions := &mqcloudv1.CreateUserOptions{
				ServiceInstanceGuid: core.StringPtr(config["SERVICE_INSTANCE_DEPLOYMENT_GUID"]),
				Email:               core.StringPtr("testuser@ibm.com"),
				Name:                core.StringPtr("testuser"),
			}

			userDetails, response, err := mqcloudService.CreateUser(createUserOptions)
			Expect(err).To(BeNil())
			config["USER_ID"] = *userDetails.ID
			Expect(response.StatusCode).To(Equal(201))
			Expect(userDetails).ToNot(BeNil())
		})
	})

	Describe(`GetUser - Get a user for an instance`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`GetUser(getUserOptions *GetUserOptions)`, func() {
			getUserOptions := &mqcloudv1.GetUserOptions{
				ServiceInstanceGuid: core.StringPtr(config["SERVICE_INSTANCE_DEPLOYMENT_GUID"]),
				UserID:              core.StringPtr(config["USER_ID"]),
			}

			userDetails, response, err := mqcloudService.GetUser(getUserOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(userDetails).ToNot(BeNil())
		})
	})

	Describe(`ListApplications - Get a list of applications for an instance`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`ListApplications(listApplicationsOptions *ListApplicationsOptions) with pagination`, func() {
			listApplicationsOptions := &mqcloudv1.ListApplicationsOptions{
				ServiceInstanceGuid: core.StringPtr(config["SERVICE_INSTANCE_DEPLOYMENT_GUID"]),
				Offset:              core.Int64Ptr(int64(0)),
				Limit:               core.Int64Ptr(int64(10)),
			}

			listApplicationsOptions.Offset = nil
			listApplicationsOptions.Limit = core.Int64Ptr(1)

			var allResults []mqcloudv1.ApplicationDetails
			for {
				applicationDetailsCollection, response, err := mqcloudService.ListApplications(listApplicationsOptions)
				Expect(err).To(BeNil())
				Expect(response.StatusCode).To(Equal(200))
				Expect(applicationDetailsCollection).ToNot(BeNil())
				allResults = append(allResults, applicationDetailsCollection.Applications...)

				listApplicationsOptions.Offset, err = applicationDetailsCollection.GetNextOffset()
				Expect(err).To(BeNil())

				if listApplicationsOptions.Offset == nil {
					break
				}
			}
			fmt.Fprintf(GinkgoWriter, "Retrieved a total of %d item(s) with pagination.\n", len(allResults))
		})
		Describe(`CreateApplication - Add an application to an instance`, func() {
			BeforeEach(func() {
				shouldSkipTest()
			})
			It(`CreateApplication(createApplicationOptions *CreateApplicationOptions)`, func() {
				createApplicationOptions := &mqcloudv1.CreateApplicationOptions{
					ServiceInstanceGuid: core.StringPtr(config["SERVICE_INSTANCE_DEPLOYMENT_GUID"]),
					Name:                core.StringPtr("test-app"),
				}

				applicationCreated, response, err := mqcloudService.CreateApplication(createApplicationOptions)
				Expect(err).To(BeNil())
				config["APPLICATION_ID"] = *applicationCreated.ID
				Expect(response.StatusCode).To(Equal(201))
				Expect(applicationCreated).ToNot(BeNil())
			})
		})

		It(`ListApplications(listApplicationsOptions *ListApplicationsOptions) using ApplicationsPager`, func() {
			listApplicationsOptions := &mqcloudv1.ListApplicationsOptions{
				ServiceInstanceGuid: core.StringPtr(config["SERVICE_INSTANCE_DEPLOYMENT_GUID"]),
				Limit:               core.Int64Ptr(int64(10)),
			}

			// Test GetNext().
			pager, err := mqcloudService.NewApplicationsPager(listApplicationsOptions)
			Expect(err).To(BeNil())
			Expect(pager).ToNot(BeNil())

			var allResults []mqcloudv1.ApplicationDetails
			for pager.HasNext() {
				nextPage, err := pager.GetNext()
				Expect(err).To(BeNil())
				Expect(nextPage).ToNot(BeNil())
				allResults = append(allResults, nextPage...)
			}

			// Test GetAll().
			pager, err = mqcloudService.NewApplicationsPager(listApplicationsOptions)
			Expect(err).To(BeNil())
			Expect(pager).ToNot(BeNil())

			allItems, err := pager.GetAll()
			Expect(err).To(BeNil())
			Expect(allItems).ToNot(BeNil())

			Expect(len(allItems)).To(Equal(len(allResults)))
			fmt.Fprintf(GinkgoWriter, "ListApplications() returned a total of %d item(s) using ApplicationsPager.\n", len(allResults))
		})
	})

	Describe(`GetApplication - Get an application for an instance`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`GetApplication(getApplicationOptions *GetApplicationOptions)`, func() {
			getApplicationOptions := &mqcloudv1.GetApplicationOptions{
				ServiceInstanceGuid: core.StringPtr(config["SERVICE_INSTANCE_DEPLOYMENT_GUID"]),
				ApplicationID:       core.StringPtr(config["APPLICATION_ID"]),
			}

			applicationDetails, response, err := mqcloudService.GetApplication(getApplicationOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(applicationDetails).ToNot(BeNil())
		})
	})

	Describe(`CreateApplicationApikey - Create a new apikey for an application`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`CreateApplicationApikey(createApplicationApikeyOptions *CreateApplicationApikeyOptions)`, func() {
			createApplicationApikeyOptions := &mqcloudv1.CreateApplicationApikeyOptions{
				ServiceInstanceGuid: core.StringPtr(config["SERVICE_INSTANCE_DEPLOYMENT_GUID"]),
				ApplicationID:       core.StringPtr(config["APPLICATION_ID"]),
				Name:                core.StringPtr("test-api-key"),
			}

			applicationApiKeyCreated, response, err := mqcloudService.CreateApplicationApikey(createApplicationApikeyOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(201))
			Expect(applicationApiKeyCreated).ToNot(BeNil())
		})
	})

	Describe(`CreateTrustStorePemCertificate - Upload a trust store certificate`, func() {
		BeforeEach(func() {
			shouldSkipTest()
			SkipTestIfQmIsNotRunning(config["QUEUE_MANAGER_ID"], mqcloudService, config["SERVICE_INSTANCE_DEPLOYMENT_GUID"])

		})
		It(`CreateTrustStorePemCertificate(createTrustStorePemCertificateOptions *CreateTrustStorePemCertificateOptions)`, func() {
			file, err := os.Open(config["TRUST_STORE_FILE_PATH"])
			if err != nil {
				fmt.Fprintf(GinkgoWriter, "Error opening file: %s \n", err.Error())
				return
			}
			defer file.Close()
			createTrustStorePemCertificateOptions := &mqcloudv1.CreateTrustStorePemCertificateOptions{
				ServiceInstanceGuid: core.StringPtr(config["SERVICE_INSTANCE_DEPLOYMENT_GUID"]),
				QueueManagerID:      core.StringPtr(config["QUEUE_MANAGER_ID"]),
				Label:               core.StringPtr("ittruststore"),
				CertificateFile:     file,
			}

			trustStoreCertificateDetails, response, err := mqcloudService.CreateTrustStorePemCertificate(createTrustStorePemCertificateOptions)
			Expect(err).To(BeNil())
			config["TRUST_STORE_CERTIFICATE_ID"] = *trustStoreCertificateDetails.ID
			time.Sleep(60 * time.Second)
			Expect(response.StatusCode).To(Equal(201))
			Expect(trustStoreCertificateDetails).ToNot(BeNil())
		})
	})

	Describe(`ListTrustStoreCertificates - List trust store certificates`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`ListTrustStoreCertificates(listTrustStoreCertificatesOptions *ListTrustStoreCertificatesOptions)`, func() {
			listTrustStoreCertificatesOptions := &mqcloudv1.ListTrustStoreCertificatesOptions{
				ServiceInstanceGuid: core.StringPtr(config["SERVICE_INSTANCE_DEPLOYMENT_GUID"]),
				QueueManagerID:      core.StringPtr(config["QUEUE_MANAGER_ID"]),
			}

			trustStoreCertificateDetailsCollection, response, err := mqcloudService.ListTrustStoreCertificates(listTrustStoreCertificatesOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(trustStoreCertificateDetailsCollection).ToNot(BeNil())
		})
	})

	Describe(`GetTrustStoreCertificate - Get a trust store certificate`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`GetTrustStoreCertificate(getTrustStoreCertificateOptions *GetTrustStoreCertificateOptions)`, func() {
			getTrustStoreCertificateOptions := &mqcloudv1.GetTrustStoreCertificateOptions{
				ServiceInstanceGuid: core.StringPtr(config["SERVICE_INSTANCE_DEPLOYMENT_GUID"]),
				QueueManagerID:      core.StringPtr(config["QUEUE_MANAGER_ID"]),
				CertificateID:       core.StringPtr(config["TRUST_STORE_CERTIFICATE_ID"]),
			}

			trustStoreCertificateDetails, response, err := mqcloudService.GetTrustStoreCertificate(getTrustStoreCertificateOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(trustStoreCertificateDetails).ToNot(BeNil())
		})
	})

	Describe(`DownloadTrustStoreCertificate - Download a queue manager's certificate from its trust store`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`DownloadTrustStoreCertificate(downloadTrustStoreCertificateOptions *DownloadTrustStoreCertificateOptions)`, func() {
			downloadTrustStoreCertificateOptions := &mqcloudv1.DownloadTrustStoreCertificateOptions{
				ServiceInstanceGuid: core.StringPtr(config["SERVICE_INSTANCE_DEPLOYMENT_GUID"]),
				QueueManagerID:      core.StringPtr(config["QUEUE_MANAGER_ID"]),
				CertificateID:       core.StringPtr(config["TRUST_STORE_CERTIFICATE_ID"]),
			}

			result, response, err := mqcloudService.DownloadTrustStoreCertificate(downloadTrustStoreCertificateOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(result).ToNot(BeNil())
		})
	})

	Describe(`CreateKeyStorePemCertificate - Upload a key store certificate`, func() {
		BeforeEach(func() {
			shouldSkipTest()
			SkipTestIfQmIsNotRunning(config["QUEUE_MANAGER_ID"], mqcloudService, config["SERVICE_INSTANCE_DEPLOYMENT_GUID"])
		})
		It(`CreateKeyStorePemCertificate(createKeyStorePemCertificateOptions *CreateKeyStorePemCertificateOptions)`, func() {
			file, err := os.Open(config["KEY_STORE_FILE_PATH"])
			if err != nil {
				fmt.Fprintf(GinkgoWriter, "Error opening file: %s \n", err.Error())
				return
			}
			defer file.Close()
			createKeyStorePemCertificateOptions := &mqcloudv1.CreateKeyStorePemCertificateOptions{
				ServiceInstanceGuid: core.StringPtr(config["SERVICE_INSTANCE_DEPLOYMENT_GUID"]),
				QueueManagerID:      core.StringPtr(config["QUEUE_MANAGER_ID"]),
				Label:               core.StringPtr("itkeystore"),
				CertificateFile:     file,
			}

			keyStoreCertificateDetails, response, err := mqcloudService.CreateKeyStorePemCertificate(createKeyStorePemCertificateOptions)
			Expect(err).To(BeNil())
			config["KEY_STORE_CERTIFICATE_ID"] = *keyStoreCertificateDetails.ID
			time.Sleep(60 * time.Second)
			Expect(response.StatusCode).To(Equal(201))
			Expect(keyStoreCertificateDetails).ToNot(BeNil())
		})
	})

	Describe(`ListKeyStoreCertificates - List key store certificates`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`ListKeyStoreCertificates(listKeyStoreCertificatesOptions *ListKeyStoreCertificatesOptions)`, func() {
			listKeyStoreCertificatesOptions := &mqcloudv1.ListKeyStoreCertificatesOptions{
				ServiceInstanceGuid: core.StringPtr(config["SERVICE_INSTANCE_DEPLOYMENT_GUID"]),
				QueueManagerID:      core.StringPtr(config["QUEUE_MANAGER_ID"]),
			}

			keyStoreCertificateDetailsCollection, response, err := mqcloudService.ListKeyStoreCertificates(listKeyStoreCertificatesOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(keyStoreCertificateDetailsCollection).ToNot(BeNil())
		})
	})

	Describe(`GetKeyStoreCertificate - Get a key store certificate for queue manager`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`GetKeyStoreCertificate(getKeyStoreCertificateOptions *GetKeyStoreCertificateOptions)`, func() {
			getKeyStoreCertificateOptions := &mqcloudv1.GetKeyStoreCertificateOptions{
				ServiceInstanceGuid: core.StringPtr(config["SERVICE_INSTANCE_DEPLOYMENT_GUID"]),
				QueueManagerID:      core.StringPtr(config["QUEUE_MANAGER_ID"]),
				CertificateID:       core.StringPtr(config["KEY_STORE_CERTIFICATE_ID"]),
			}

			keyStoreCertificateDetails, response, err := mqcloudService.GetKeyStoreCertificate(getKeyStoreCertificateOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(keyStoreCertificateDetails).ToNot(BeNil())
		})
	})

	Describe(`DownloadKeyStoreCertificate - Download a queue manager's certificate from its key store`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`DownloadKeyStoreCertificate(downloadKeyStoreCertificateOptions *DownloadKeyStoreCertificateOptions)`, func() {
			downloadKeyStoreCertificateOptions := &mqcloudv1.DownloadKeyStoreCertificateOptions{
				ServiceInstanceGuid: core.StringPtr(config["SERVICE_INSTANCE_DEPLOYMENT_GUID"]),
				QueueManagerID:      core.StringPtr(config["QUEUE_MANAGER_ID"]),
				CertificateID:       core.StringPtr(config["KEY_STORE_CERTIFICATE_ID"]),
			}

			result, response, err := mqcloudService.DownloadKeyStoreCertificate(downloadKeyStoreCertificateOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(result).ToNot(BeNil())
		})
	})

	Describe(`GetCertificateAmsChannels - Get the AMS channels that are configured with this key store certificate`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`GetCertificateAmsChannels(getCertificateAmsChannelsOptions *GetCertificateAmsChannelsOptions)`, func() {
			getCertificateAmsChannelsOptions := &mqcloudv1.GetCertificateAmsChannelsOptions{
				QueueManagerID:      core.StringPtr(config["QUEUE_MANAGER_ID"]),
				CertificateID:       core.StringPtr(config["KEY_STORE_CERTIFICATE_ID"]),
				ServiceInstanceGuid: core.StringPtr(config["SERVICE_INSTANCE_DEPLOYMENT_GUID"]),
			}

			channelsDetails, response, err := mqcloudService.GetCertificateAmsChannels(getCertificateAmsChannelsOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(channelsDetails).ToNot(BeNil())
		})
	})

	Describe(`SetCertificateAmsChannels - Update the AMS channels that are configured with this key store certificate`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`SetCertificateAmsChannels(setCertificateAmsChannelsOptions *SetCertificateAmsChannelsOptions)`, func() {
			channelDetailsModel := &mqcloudv1.ChannelDetails{
				Name: core.StringPtr("CLOUD.APP.SVRCONN"),
			}

			setCertificateAmsChannelsOptions := &mqcloudv1.SetCertificateAmsChannelsOptions{
				QueueManagerID:      core.StringPtr(config["QUEUE_MANAGER_ID"]),
				CertificateID:       core.StringPtr(config["KEY_STORE_CERTIFICATE_ID"]),
				ServiceInstanceGuid: core.StringPtr(config["SERVICE_INSTANCE_DEPLOYMENT_GUID"]),
				Channels:            []mqcloudv1.ChannelDetails{*channelDetailsModel},
				UpdateStrategy:      core.StringPtr("append"),
			}

			channelsDetails, response, err := mqcloudService.SetCertificateAmsChannels(setCertificateAmsChannelsOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(channelsDetails).ToNot(BeNil())

			// Putting back ams configs to nil, so that we can delete the keystore certificate and also test both update strategies
			setCertificateAmsChannelsOptions.Channels = []mqcloudv1.ChannelDetails{}
			setCertificateAmsChannelsOptions.UpdateStrategy = core.StringPtr("replace")
			channelsDetails, response, err = mqcloudService.SetCertificateAmsChannels(setCertificateAmsChannelsOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(channelsDetails).ToNot(BeNil())
		})
	})

	Describe(`CreateVirtualPrivateEndpointGateway - Create a new virtual private endpoint gateway`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`CreateVirtualPrivateEndpointGateway(createVirtualPrivateEndpointGatewayOptions *CreateVirtualPrivateEndpointGatewayOptions)`, func() {
			createVirtualPrivateEndpointGatewayOptions := &mqcloudv1.CreateVirtualPrivateEndpointGatewayOptions{
				ServiceInstanceGuid: core.StringPtr(config["SERVICE_INSTANCE_CAPACITY_GUID"]),
				Name:                core.StringPtr("testvpegforit"),
				TargetCrn:           core.StringPtr(config["TARGET_CRN"]),
				TrustedProfile:      core.StringPtr(config["TRUSTED_PROFILE"]),
			}

			createVirtualPrivateEndpointGatewayOptions.TrustedProfile = nil
			virtualPrivateEndpointGatewayDetails, response, err := mqcloudService.CreateVirtualPrivateEndpointGateway(createVirtualPrivateEndpointGatewayOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(201))

			Expect(virtualPrivateEndpointGatewayDetails).ToNot(BeNil())
			gatewayid := *virtualPrivateEndpointGatewayDetails.ID
			config["VIRTUAL_PRIVATE_ENDPOINT_GATEWAY_GUID"] = *virtualPrivateEndpointGatewayDetails.ID
			Expect(gatewayid).To(Equal(config["VIRTUAL_PRIVATE_ENDPOINT_GATEWAY_GUID"]))
		})
	})

	Describe(`ListVirtualPrivateEndpointGateways - Get a list of information for all virtual private endpoint gateways`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`ListVirtualPrivateEndpointGateways(listVirtualPrivateEndpointGatewaysOptions *ListVirtualPrivateEndpointGatewaysOptions) with pagination`, func() {
			listVirtualPrivateEndpointGatewaysOptions := &mqcloudv1.ListVirtualPrivateEndpointGatewaysOptions{
				ServiceInstanceGuid: core.StringPtr(config["SERVICE_INSTANCE_CAPACITY_GUID"]),
				TrustedProfile:      core.StringPtr(config["TRUSTED_PROFILE"]),
				Start:               core.StringPtr(""),
				Limit:               core.Int64Ptr(int64(10)),
			}

			listVirtualPrivateEndpointGatewaysOptions.TrustedProfile = nil
			listVirtualPrivateEndpointGatewaysOptions.Start = nil
			listVirtualPrivateEndpointGatewaysOptions.Limit = core.Int64Ptr(1)

			var allResults []mqcloudv1.VirtualPrivateEndpointGatewayDetails
			for {

				virtualPrivateEndpointGatewayDetailsCollection, response, err := mqcloudService.ListVirtualPrivateEndpointGateways(listVirtualPrivateEndpointGatewaysOptions)

				Expect(err).To(BeNil())
				Expect(response.StatusCode).To(Equal(200))
				Expect(virtualPrivateEndpointGatewayDetailsCollection).ToNot(BeNil())
				allResults = append(allResults, virtualPrivateEndpointGatewayDetailsCollection.VirtualPrivateEndpointGateways...)

				listVirtualPrivateEndpointGatewaysOptions.Start, err = virtualPrivateEndpointGatewayDetailsCollection.GetNextStart()
				Expect(err).To(BeNil())
				if listVirtualPrivateEndpointGatewaysOptions.Start == nil {

					break
				}
				fmt.Println("START=" + *listVirtualPrivateEndpointGatewaysOptions.Start)
			}
			fmt.Fprintf(GinkgoWriter, "Retrieved a total of %d item(s) with pagination.\n", len(allResults))
		})
		It(`ListVirtualPrivateEndpointGateways(listVirtualPrivateEndpointGatewaysOptions *ListVirtualPrivateEndpointGatewaysOptions) using VirtualPrivateEndpointGatewaysPager`, func() {
			listVirtualPrivateEndpointGatewaysOptions := &mqcloudv1.ListVirtualPrivateEndpointGatewaysOptions{
				ServiceInstanceGuid: core.StringPtr(config["SERVICE_INSTANCE_CAPACITY_GUID"]),
				TrustedProfile:      core.StringPtr(config["TRUSTED_PROFILE"]),
				Limit:               core.Int64Ptr(int64(10)),
			}

			// Test GetNext().
			listVirtualPrivateEndpointGatewaysOptions.TrustedProfile = nil
			pager, err := mqcloudService.NewVirtualPrivateEndpointGatewaysPager(listVirtualPrivateEndpointGatewaysOptions)
			Expect(err).To(BeNil())
			Expect(pager).ToNot(BeNil())

			var allResults []mqcloudv1.VirtualPrivateEndpointGatewayDetails
			for pager.HasNext() {
				nextPage, err := pager.GetNext()
				Expect(err).To(BeNil())
				Expect(nextPage).ToNot(BeNil())
				allResults = append(allResults, nextPage...)
			}

			// Test GetAll().
			pager, err = mqcloudService.NewVirtualPrivateEndpointGatewaysPager(listVirtualPrivateEndpointGatewaysOptions)
			Expect(err).To(BeNil())
			Expect(pager).ToNot(BeNil())

			allItems, err := pager.GetAll()
			Expect(err).To(BeNil())
			//		Expect(allItems).ToNot(BeNil())

			Expect(len(allItems)).To(Equal(len(allResults)))
			fmt.Fprintf(GinkgoWriter, "ListVirtualPrivateEndpointGateways() returned a total of %d item(s) using VirtualPrivateEndpointGatewaysPager.\n", len(allResults))
		})
	})

	Describe(`GetVirtualPrivateEndpointGateway - Display the information for a specific virtual private endpoint gateway`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`GetVirtualPrivateEndpointGateway(getVirtualPrivateEndpointGatewayOptions *GetVirtualPrivateEndpointGatewayOptions)`, func() {
			getVirtualPrivateEndpointGatewayOptions := &mqcloudv1.GetVirtualPrivateEndpointGatewayOptions{
				ServiceInstanceGuid:               core.StringPtr(config["SERVICE_INSTANCE_CAPACITY_GUID"]),
				VirtualPrivateEndpointGatewayGuid: core.StringPtr(config["VIRTUAL_PRIVATE_ENDPOINT_GATEWAY_GUID"]),
				TrustedProfile:                    core.StringPtr(config["TRUSTED_PROFILE"]),
			}

			getVirtualPrivateEndpointGatewayOptions.TrustedProfile = nil
			virtualPrivateEndpointGatewayDetails, response, err := mqcloudService.GetVirtualPrivateEndpointGateway(getVirtualPrivateEndpointGatewayOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(virtualPrivateEndpointGatewayDetails).ToNot(BeNil())
		})
	})

	Describe(`DeleteUser - Delete a user for an instance`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`DeleteUser(deleteUserOptions *DeleteUserOptions)`, func() {
			deleteUserOptions := &mqcloudv1.DeleteUserOptions{
				ServiceInstanceGuid: core.StringPtr(config["SERVICE_INSTANCE_DEPLOYMENT_GUID"]),
				UserID:              core.StringPtr(config["USER_ID"]),
			}

			response, err := mqcloudService.DeleteUser(deleteUserOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(204))
		})
	})

	Describe(`DeleteApplication - Delete an application from an instance`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`DeleteApplication(deleteApplicationOptions *DeleteApplicationOptions)`, func() {
			deleteApplicationOptions := &mqcloudv1.DeleteApplicationOptions{
				ServiceInstanceGuid: core.StringPtr(config["SERVICE_INSTANCE_DEPLOYMENT_GUID"]),
				ApplicationID:       core.StringPtr(config["APPLICATION_ID"]),
			}

			response, err := mqcloudService.DeleteApplication(deleteApplicationOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(204))
		})
	})

	Describe(`DeleteTrustStoreCertificate - Delete a trust store certificate`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`DeleteTrustStoreCertificate(deleteTrustStoreCertificateOptions *DeleteTrustStoreCertificateOptions)`, func() {
			deleteTrustStoreCertificateOptions := &mqcloudv1.DeleteTrustStoreCertificateOptions{
				ServiceInstanceGuid: core.StringPtr(config["SERVICE_INSTANCE_DEPLOYMENT_GUID"]),
				QueueManagerID:      core.StringPtr(config["QUEUE_MANAGER_ID"]),
				CertificateID:       core.StringPtr(config["TRUST_STORE_CERTIFICATE_ID"]),
			}

			response, err := mqcloudService.DeleteTrustStoreCertificate(deleteTrustStoreCertificateOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(204))
		})
	})

	Describe(`DeleteKeyStoreCertificate - Delete a queue manager's key store certificate`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`DeleteKeyStoreCertificate(deleteKeyStoreCertificateOptions *DeleteKeyStoreCertificateOptions)`, func() {
			deleteKeyStoreCertificateOptions := &mqcloudv1.DeleteKeyStoreCertificateOptions{
				ServiceInstanceGuid: core.StringPtr(config["SERVICE_INSTANCE_DEPLOYMENT_GUID"]),
				QueueManagerID:      core.StringPtr(config["QUEUE_MANAGER_ID"]),
				CertificateID:       core.StringPtr(config["KEY_STORE_CERTIFICATE_ID"]),
			}

			response, err := mqcloudService.DeleteKeyStoreCertificate(deleteKeyStoreCertificateOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(204))
		})
	})

	Describe(`DeleteQueueManager - Delete a queue manager`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`DeleteQueueManager(deleteQueueManagerOptions *DeleteQueueManagerOptions)`, func() {
			deleteQueueManagerOptions := &mqcloudv1.DeleteQueueManagerOptions{
				ServiceInstanceGuid: core.StringPtr(config["SERVICE_INSTANCE_DEPLOYMENT_GUID"]),
				QueueManagerID:      core.StringPtr(config["QUEUE_MANAGER_ID"]),
			}

			queueManagerTaskStatus, response, err := mqcloudService.DeleteQueueManager(deleteQueueManagerOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(202))
			Expect(queueManagerTaskStatus).ToNot(BeNil())
		})
	})

	Describe(`DeleteVirtualPrivateEndpointGateway - Delete a specific virtual private endpoint gateway`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`DeleteVirtualPrivateEndpointGateway(deleteVirtualPrivateEndpointGatewayOptions *DeleteVirtualPrivateEndpointGatewayOptions)`, func() {
			deleteVirtualPrivateEndpointGatewayOptions := &mqcloudv1.DeleteVirtualPrivateEndpointGatewayOptions{
				ServiceInstanceGuid:               core.StringPtr(config["SERVICE_INSTANCE_CAPACITY_GUID"]),
				VirtualPrivateEndpointGatewayGuid: core.StringPtr(config["VIRTUAL_PRIVATE_ENDPOINT_GATEWAY_GUID"]),
				TrustedProfile:                    core.StringPtr(config["TRUSTED_PROFILE"]),
			}
			deleteVirtualPrivateEndpointGatewayOptions.TrustedProfile = nil
			response, err := mqcloudService.DeleteVirtualPrivateEndpointGateway(deleteVirtualPrivateEndpointGatewayOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(204))
		})
	})
})

//
// Utility functions are declared in the unit test file
//
