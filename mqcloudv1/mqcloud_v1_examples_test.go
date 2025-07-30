//go:build examples

/**
 * (C) Copyright IBM Corp. 2025.
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
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/mqcloud-go-sdk/mqcloudv1"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

// This file provides an example of how to use the mqcloud service.
//
// The following configuration properties are assumed to be defined:
// MQCLOUD_URL=<service base url>
// MQCLOUD_AUTH_TYPE=iam
// MQCLOUD_APIKEY=<IAM apikey>
// MQCLOUD_AUTH_URL=<IAM token service base URL - omit this if using the production environment>
//
// These configuration properties can be exported as environment variables, or stored
// in a configuration file and then:
// export IBM_CREDENTIALS_FILE=<name of configuration file>
var _ = Describe(`MqcloudV1 Examples Tests`, func() {

	const externalConfigFile = "../mqcloud_v1.env"

	var (
		mqcloudService *mqcloudv1.MqcloudV1
		config         map[string]string
	)

	var shouldSkipTest = func() {
		Skip("External configuration is not available, skipping examples...")
	}

	Describe(`External configuration`, func() {
		It("Successfully load the configuration", func() {
			var err error
			_, err = os.Stat(externalConfigFile)
			if err != nil {
				Skip("External configuration file not found, skipping examples: " + err.Error())
			}

			os.Setenv("IBM_CREDENTIALS_FILE", externalConfigFile)
			config, err = core.GetServiceProperties(mqcloudv1.DefaultServiceName)
			if err != nil {
				Skip("Error loading service properties, skipping examples: " + err.Error())
			} else if len(config) == 0 {
				Skip("Unable to load service properties, skipping examples")
			}

			shouldSkipTest = func() {}
		})
	})

	Describe(`Client initialization`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It("Successfully construct the service client instance", func() {
			var err error

			// begin-common

			mqcloudServiceOptions := &mqcloudv1.MqcloudV1Options{}

			mqcloudService, err = mqcloudv1.NewMqcloudV1UsingExternalConfig(mqcloudServiceOptions)

			if err != nil {
				panic(err)
			}

			// end-common

			Expect(mqcloudService).ToNot(BeNil())
		})
	})

	Describe(`MqcloudV1 request examples`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`GetUsageDetails request example`, func() {
			fmt.Println("\nGetUsageDetails() result:")
			// begin-get_usage_details

			getUsageDetailsOptions := mqcloudService.NewGetUsageDetailsOptions(
				config["SERVICE_INSTANCE_DEPLOYMENT_GUID"],
			)

			usage, response, err := mqcloudService.GetUsageDetails(getUsageDetailsOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(usage, "", "  ")
			fmt.Println(string(b))

			// end-get_usage_details

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(usage).ToNot(BeNil())
		})
		It(`GetOptions request example`, func() {
			fmt.Println("\nGetOptions() result:")
			// begin-get_options

			getOptionsOptions := mqcloudService.NewGetOptionsOptions(
				config["SERVICE_INSTANCE_DEPLOYMENT_GUID"],
			)

			configurationOptions, response, err := mqcloudService.GetOptions(getOptionsOptions)
			if err != nil {
				panic(err)
			}
			config["LOCATION"] = configurationOptions.Locations[0]
			b, _ := json.MarshalIndent(configurationOptions, "", "  ")
			fmt.Println(string(b))

			// end-get_options

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(configurationOptions).ToNot(BeNil())
		})
		It(`CreateQueueManager request example`, func() {
			fmt.Println("\nCreateQueueManager() result:")
			// begin-create_queue_manager

			createQueueManagerOptions := mqcloudService.NewCreateQueueManagerOptions(
				config["SERVICE_INSTANCE_DEPLOYMENT_GUID"],
				"testqm",
				config["LOCATION"],
				"small",
			)
			createQueueManagerOptions.SetVersion(config["VERSION"])

			queueManagerTaskStatus, response, err := mqcloudService.CreateQueueManager(createQueueManagerOptions)
			if err != nil {
				panic(err)
			}
			config["QUEUE_MANAGER_ID"] = *queueManagerTaskStatus.QueueManagerID
			b, _ := json.MarshalIndent(queueManagerTaskStatus, "", "  ")
			fmt.Println(string(b))

			// end-create_queue_manager

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(202))
			Expect(queueManagerTaskStatus).ToNot(BeNil())
		})
		It(`ListQueueManagers request example`, func() {
			fmt.Println("\nListQueueManagers() result:")
			// begin-list_queue_managers
			listQueueManagersOptions := &mqcloudv1.ListQueueManagersOptions{
				ServiceInstanceGuid: core.StringPtr(config["SERVICE_INSTANCE_DEPLOYMENT_GUID"]),
				Limit:               core.Int64Ptr(int64(10)),
			}

			pager, err := mqcloudService.NewQueueManagersPager(listQueueManagersOptions)
			if err != nil {
				panic(err)
			}

			var allResults []mqcloudv1.QueueManagerDetails
			for pager.HasNext() {
				nextPage, err := pager.GetNext()
				if err != nil {
					panic(err)
				}
				allResults = append(allResults, nextPage...)
			}
			b, _ := json.MarshalIndent(allResults, "", "  ")
			fmt.Println(string(b))
			// end-list_queue_managers
		})
		It(`GetQueueManager request example`, func() {
			fmt.Println("\nGetQueueManager() result:")
			// begin-get_queue_manager

			getQueueManagerOptions := mqcloudService.NewGetQueueManagerOptions(
				config["SERVICE_INSTANCE_DEPLOYMENT_GUID"],
				config["QUEUE_MANAGER_ID"],
			)

			queueManagerDetails, response, err := mqcloudService.GetQueueManager(getQueueManagerOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(queueManagerDetails, "", "  ")
			fmt.Println(string(b))

			// end-get_queue_manager

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(queueManagerDetails).ToNot(BeNil())
		})
		It(`SetQueueManagerVersion request example`, func() {

			SkipTestIfQmIsNotRunning(config["QUEUE_MANAGER_ID"], mqcloudService, config["SERVICE_INSTANCE_DEPLOYMENT_GUID"])
			fmt.Println("\nSetQueueManagerVersion() result:")
			// begin-set_queue_manager_version

			setQueueManagerVersionOptions := mqcloudService.NewSetQueueManagerVersionOptions(
				config["SERVICE_INSTANCE_DEPLOYMENT_GUID"],
				config["QUEUE_MANAGER_ID"],
				config["VERSION_UPGRADE"],
			)

			queueManagerTaskStatus, response, err := mqcloudService.SetQueueManagerVersion(setQueueManagerVersionOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(queueManagerTaskStatus, "", "  ")
			fmt.Println(string(b))

			// end-set_queue_manager_version

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(202))
			Expect(queueManagerTaskStatus).ToNot(BeNil())
		})
		It(`GetQueueManagerAvailableUpgradeVersions request example`, func() {
			fmt.Println("\nGetQueueManagerAvailableUpgradeVersions() result:")
			// begin-get_queue_manager_available_upgrade_versions

			getQueueManagerAvailableUpgradeVersionsOptions := mqcloudService.NewGetQueueManagerAvailableUpgradeVersionsOptions(
				config["SERVICE_INSTANCE_DEPLOYMENT_GUID"],
				config["QUEUE_MANAGER_ID"],
			)

			queueManagerVersionUpgrades, response, err := mqcloudService.GetQueueManagerAvailableUpgradeVersions(getQueueManagerAvailableUpgradeVersionsOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(queueManagerVersionUpgrades, "", "  ")
			fmt.Println(string(b))

			// end-get_queue_manager_available_upgrade_versions

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(queueManagerVersionUpgrades).ToNot(BeNil())
		})
		It(`GetQueueManagerConnectionInfo request example`, func() {
			fmt.Println("\nGetQueueManagerConnectionInfo() result:")
			// begin-get_queue_manager_connection_info

			getQueueManagerConnectionInfoOptions := mqcloudService.NewGetQueueManagerConnectionInfoOptions(
				config["SERVICE_INSTANCE_DEPLOYMENT_GUID"],
				config["QUEUE_MANAGER_ID"],
			)

			connectionInfo, response, err := mqcloudService.GetQueueManagerConnectionInfo(getQueueManagerConnectionInfoOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(connectionInfo, "", "  ")
			fmt.Println(string(b))

			// end-get_queue_manager_connection_info

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(connectionInfo).ToNot(BeNil())
		})
		It(`GetQueueManagerStatus request example`, func() {
			fmt.Println("\nGetQueueManagerStatus() result:")
			// begin-get_queue_manager_status

			getQueueManagerStatusOptions := mqcloudService.NewGetQueueManagerStatusOptions(
				config["SERVICE_INSTANCE_DEPLOYMENT_GUID"],
				config["QUEUE_MANAGER_ID"],
			)

			queueManagerStatus, response, err := mqcloudService.GetQueueManagerStatus(getQueueManagerStatusOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(queueManagerStatus, "", "  ")
			fmt.Println(string(b))

			// end-get_queue_manager_status

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(queueManagerStatus).ToNot(BeNil())
		})
		It(`ListUsers request example`, func() {
			fmt.Println("\nListUsers() result:")
			// begin-list_users
			listUsersOptions := &mqcloudv1.ListUsersOptions{
				ServiceInstanceGuid: core.StringPtr(config["SERVICE_INSTANCE_DEPLOYMENT_GUID"]),
				Limit:               core.Int64Ptr(int64(10)),
			}

			pager, err := mqcloudService.NewUsersPager(listUsersOptions)
			if err != nil {
				panic(err)
			}

			var allResults []mqcloudv1.UserDetails
			for pager.HasNext() {
				nextPage, err := pager.GetNext()
				if err != nil {
					panic(err)
				}
				allResults = append(allResults, nextPage...)
			}
			b, _ := json.MarshalIndent(allResults, "", "  ")
			fmt.Println(string(b))
			// end-list_users
		})
		It(`CreateUser request example`, func() {
			fmt.Println("\nCreateUser() result:")
			// begin-create_user

			createUserOptions := mqcloudService.NewCreateUserOptions(
				config["SERVICE_INSTANCE_DEPLOYMENT_GUID"],
				"testuser@ibm.com",
				"testuser",
			)

			userDetails, response, err := mqcloudService.CreateUser(createUserOptions)
			if err != nil {
				panic(err)
			}
			config["USER_ID"] = *userDetails.ID
			b, _ := json.MarshalIndent(userDetails, "", "  ")
			fmt.Println(string(b))

			// end-create_user

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(201))
			Expect(userDetails).ToNot(BeNil())
		})
		It(`GetUser request example`, func() {
			fmt.Println("\nGetUser() result:")
			// begin-get_user

			getUserOptions := mqcloudService.NewGetUserOptions(
				config["SERVICE_INSTANCE_DEPLOYMENT_GUID"],
				config["USER_ID"],
			)

			userDetails, response, err := mqcloudService.GetUser(getUserOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(userDetails, "", "  ")
			fmt.Println(string(b))

			// end-get_user

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(userDetails).ToNot(BeNil())
		})
		It(`SetUserName request example`, func() {
			fmt.Println("\nSetUserName() result:")
			// begin-set_user_name

			setUserNameOptions := mqcloudService.NewSetUserNameOptions(
				config["SERVICE_INSTANCE_DEPLOYMENT_GUID"],
				config["USER_ID"],
				"testString",
			)

			userDetails, response, err := mqcloudService.SetUserName(setUserNameOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(userDetails, "", "  ")
			fmt.Println(string(b))

			// end-set_user_name

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(userDetails).ToNot(BeNil())
		})
		It(`ListApplications request example`, func() {
			fmt.Println("\nListApplications() result:")
			// begin-list_applications
			listApplicationsOptions := &mqcloudv1.ListApplicationsOptions{
				ServiceInstanceGuid: core.StringPtr(config["SERVICE_INSTANCE_DEPLOYMENT_GUID"]),
				Limit:               core.Int64Ptr(int64(10)),
			}

			pager, err := mqcloudService.NewApplicationsPager(listApplicationsOptions)
			if err != nil {
				panic(err)
			}

			var allResults []mqcloudv1.ApplicationDetails
			for pager.HasNext() {
				nextPage, err := pager.GetNext()
				if err != nil {
					panic(err)
				}
				allResults = append(allResults, nextPage...)
			}
			b, _ := json.MarshalIndent(allResults, "", "  ")
			fmt.Println(string(b))
			// end-list_applications
		})
		It(`CreateApplication request example`, func() {
			fmt.Println("\nCreateApplication() result:")
			// begin-create_application

			createApplicationOptions := mqcloudService.NewCreateApplicationOptions(
				config["SERVICE_INSTANCE_DEPLOYMENT_GUID"],
				"test-app",
			)

			applicationCreated, response, err := mqcloudService.CreateApplication(createApplicationOptions)
			if err != nil {
				panic(err)
			}
			config["APPLICATION_ID"] = *applicationCreated.ID
			b, _ := json.MarshalIndent(applicationCreated, "", "  ")
			fmt.Println(string(b))

			// end-create_application

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(201))
			Expect(applicationCreated).ToNot(BeNil())
		})
		It(`GetApplication request example`, func() {
			fmt.Println("\nGetApplication() result:")
			// begin-get_application

			getApplicationOptions := mqcloudService.NewGetApplicationOptions(
				config["SERVICE_INSTANCE_DEPLOYMENT_GUID"],
				config["APPLICATION_ID"],
			)

			applicationDetails, response, err := mqcloudService.GetApplication(getApplicationOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(applicationDetails, "", "  ")
			fmt.Println(string(b))

			// end-get_application

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(applicationDetails).ToNot(BeNil())
		})
		It(`SetApplicationName request example`, func() {
			fmt.Println("\nSetApplicationName() result:")
			// begin-set_application_name

			setApplicationNameOptions := mqcloudService.NewSetApplicationNameOptions(
				config["SERVICE_INSTANCE_DEPLOYMENT_GUID"],
				config["APPLICATION_ID"],
				"testString",
			)

			applicationDetails, response, err := mqcloudService.SetApplicationName(setApplicationNameOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(applicationDetails, "", "  ")
			fmt.Println(string(b))

			// end-set_application_name

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(applicationDetails).ToNot(BeNil())
		})
		It(`CreateApplicationApikey request example`, func() {
			fmt.Println("\nCreateApplicationApikey() result:")
			// begin-create_application_apikey

			createApplicationApikeyOptions := mqcloudService.NewCreateApplicationApikeyOptions(
				config["SERVICE_INSTANCE_DEPLOYMENT_GUID"],
				config["APPLICATION_ID"],
				"test-api-key",
			)

			applicationApiKeyCreated, response, err := mqcloudService.CreateApplicationApikey(createApplicationApikeyOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(applicationApiKeyCreated, "", "  ")
			fmt.Println(string(b))

			// end-create_application_apikey

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(201))
			Expect(applicationApiKeyCreated).ToNot(BeNil())
		})
		It(`CreateTrustStorePemCertificate request example`, func() {

			SkipTestIfQmIsNotRunning(config["QUEUE_MANAGER_ID"], mqcloudService, config["SERVICE_INSTANCE_DEPLOYMENT_GUID"])
			fmt.Println("\nCreateTrustStorePemCertificate() result:")
			// begin-create_trust_store_pem_certificate
			file, err := os.Open(config["TRUST_STORE_FILE_PATH"])
			if err != nil {
				fmt.Fprintf(GinkgoWriter, "Error opening file: %s \n", err.Error())
				return
			}
			defer file.Close()

			createTrustStorePemCertificateOptions := mqcloudService.NewCreateTrustStorePemCertificateOptions(
				config["SERVICE_INSTANCE_DEPLOYMENT_GUID"],
				config["QUEUE_MANAGER_ID"],
				"truststore",
				file,
			)

			trustStoreCertificateDetails, response, err := mqcloudService.CreateTrustStorePemCertificate(createTrustStorePemCertificateOptions)
			if err != nil {
				panic(err)
			}
			config["TRUST_STORE_CERTIFICATE_ID"] = *trustStoreCertificateDetails.ID
			b, _ := json.MarshalIndent(trustStoreCertificateDetails, "", "  ")
			fmt.Println(string(b))
			time.Sleep(30 * time.Second) // Add sleep for cert to upload and apply to qm and be available during delete
			// end-create_trust_store_pem_certificate

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(201))
			Expect(trustStoreCertificateDetails).ToNot(BeNil())
		})
		It(`ListTrustStoreCertificates request example`, func() {
			fmt.Println("\nListTrustStoreCertificates() result:")
			// begin-list_trust_store_certificates

			listTrustStoreCertificatesOptions := mqcloudService.NewListTrustStoreCertificatesOptions(
				config["SERVICE_INSTANCE_DEPLOYMENT_GUID"],
				config["QUEUE_MANAGER_ID"],
			)

			trustStoreCertificateDetailsCollection, response, err := mqcloudService.ListTrustStoreCertificates(listTrustStoreCertificatesOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(trustStoreCertificateDetailsCollection, "", "  ")
			fmt.Println(string(b))

			// end-list_trust_store_certificates

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(trustStoreCertificateDetailsCollection).ToNot(BeNil())
		})
		It(`GetTrustStoreCertificate request example`, func() {
			fmt.Println("\nGetTrustStoreCertificate() result:")
			// begin-get_trust_store_certificate

			getTrustStoreCertificateOptions := mqcloudService.NewGetTrustStoreCertificateOptions(
				config["SERVICE_INSTANCE_DEPLOYMENT_GUID"],
				config["QUEUE_MANAGER_ID"],
				config["TRUST_STORE_CERTIFICATE_ID"],
			)

			trustStoreCertificateDetails, response, err := mqcloudService.GetTrustStoreCertificate(getTrustStoreCertificateOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(trustStoreCertificateDetails, "", "  ")
			fmt.Println(string(b))

			// end-get_trust_store_certificate

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(trustStoreCertificateDetails).ToNot(BeNil())
		})
		It(`DownloadTrustStoreCertificate request example`, func() {
			fmt.Println("\nDownloadTrustStoreCertificate() result:")
			// begin-download_trust_store_certificate

			downloadTrustStoreCertificateOptions := mqcloudService.NewDownloadTrustStoreCertificateOptions(
				config["SERVICE_INSTANCE_DEPLOYMENT_GUID"],
				config["QUEUE_MANAGER_ID"],
				config["TRUST_STORE_CERTIFICATE_ID"],
			)

			result, response, err := mqcloudService.DownloadTrustStoreCertificate(downloadTrustStoreCertificateOptions)
			if err != nil {
				panic(err)
			}
			if result != nil {
				defer result.Close()
				outFile, err := os.Create("result.out")
				if err != nil {
					panic(err)
				}
				defer outFile.Close()
				_, err = io.Copy(outFile, result)
				if err != nil {
					panic(err)
				}
			}

			// end-download_trust_store_certificate

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(result).ToNot(BeNil())
		})
		It(`CreateKeyStorePemCertificate request example`, func() {

			SkipTestIfQmIsNotRunning(config["QUEUE_MANAGER_ID"], mqcloudService, config["SERVICE_INSTANCE_DEPLOYMENT_GUID"])
			fmt.Println("\nCreateKeyStorePemCertificate() result:")
			// begin-create_key_store_pem_certificate
			file, err := os.Open(config["KEY_STORE_FILE_PATH"])
			if err != nil {
				fmt.Fprintf(GinkgoWriter, "Error opening file: %s \n", err.Error())
				return
			}
			defer file.Close()
			createKeyStorePemCertificateOptions := mqcloudService.NewCreateKeyStorePemCertificateOptions(
				config["SERVICE_INSTANCE_DEPLOYMENT_GUID"],
				config["QUEUE_MANAGER_ID"],
				"keystore",
				file,
			)

			keyStoreCertificateDetails, response, err := mqcloudService.CreateKeyStorePemCertificate(createKeyStorePemCertificateOptions)
			if err != nil {
				panic(err)
			}
			config["KEY_STORE_CERTIFICATE_ID"] = *keyStoreCertificateDetails.ID
			b, _ := json.MarshalIndent(keyStoreCertificateDetails, "", "  ")
			fmt.Println(string(b))
			time.Sleep(30 * time.Second) // Add sleep for cert to upload and apply to qm and be available during delete
			// end-create_key_store_pem_certificate

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(201))
			Expect(keyStoreCertificateDetails).ToNot(BeNil())
		})
		It(`ListKeyStoreCertificates request example`, func() {
			fmt.Println("\nListKeyStoreCertificates() result:")
			// begin-list_key_store_certificates

			listKeyStoreCertificatesOptions := mqcloudService.NewListKeyStoreCertificatesOptions(
				config["SERVICE_INSTANCE_DEPLOYMENT_GUID"],
				config["QUEUE_MANAGER_ID"],
			)

			keyStoreCertificateDetailsCollection, response, err := mqcloudService.ListKeyStoreCertificates(listKeyStoreCertificatesOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(keyStoreCertificateDetailsCollection, "", "  ")
			fmt.Println(string(b))

			// end-list_key_store_certificates

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(keyStoreCertificateDetailsCollection).ToNot(BeNil())
		})
		It(`GetKeyStoreCertificate request example`, func() {
			fmt.Println("\nGetKeyStoreCertificate() result:")
			// begin-get_key_store_certificate

			getKeyStoreCertificateOptions := mqcloudService.NewGetKeyStoreCertificateOptions(
				config["SERVICE_INSTANCE_DEPLOYMENT_GUID"],
				config["QUEUE_MANAGER_ID"],
				config["KEY_STORE_CERTIFICATE_ID"],
			)

			keyStoreCertificateDetails, response, err := mqcloudService.GetKeyStoreCertificate(getKeyStoreCertificateOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(keyStoreCertificateDetails, "", "  ")
			fmt.Println(string(b))

			// end-get_key_store_certificate

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(keyStoreCertificateDetails).ToNot(BeNil())
		})
		It(`DownloadKeyStoreCertificate request example`, func() {
			fmt.Println("\nDownloadKeyStoreCertificate() result:")
			// begin-download_key_store_certificate

			downloadKeyStoreCertificateOptions := mqcloudService.NewDownloadKeyStoreCertificateOptions(
				config["SERVICE_INSTANCE_DEPLOYMENT_GUID"],
				config["QUEUE_MANAGER_ID"],
				config["KEY_STORE_CERTIFICATE_ID"],
			)

			result, response, err := mqcloudService.DownloadKeyStoreCertificate(downloadKeyStoreCertificateOptions)
			if err != nil {
				panic(err)
			}
			if result != nil {
				defer result.Close()
				outFile, err := os.Create("result.out")
				if err != nil {
					panic(err)
				}
				defer outFile.Close()
				_, err = io.Copy(outFile, result)
				if err != nil {
					panic(err)
				}
			}

			// end-download_key_store_certificate

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(result).ToNot(BeNil())
		})
		It(`GetCertificateAmsChannels request example`, func() {
			fmt.Println("\nGetCertificateAmsChannels() result:")
			// begin-get_certificate_ams_channels

			getCertificateAmsChannelsOptions := mqcloudService.NewGetCertificateAmsChannelsOptions(
				config["QUEUE_MANAGER_ID"],
				config["KEY_STORE_CERTIFICATE_ID"],
				config["SERVICE_INSTANCE_DEPLOYMENT_GUID"],
			)

			channelsDetails, response, err := mqcloudService.GetCertificateAmsChannels(getCertificateAmsChannelsOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(channelsDetails, "", "  ")
			fmt.Println(string(b))

			// end-get_certificate_ams_channels

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(channelsDetails).ToNot(BeNil())
		})
		It(`SetCertificateAmsChannels request example`, func() {
			fmt.Println("\nSetCertificateAmsChannels() result:")
			// begin-set_certificate_ams_channels

			// Construct an instance of the ChannelDetails model
			channelDetailsModel := new(mqcloudv1.ChannelDetails)
			channelDetailsModel.Name = core.StringPtr("testString")

			setCertificateAmsChannelsOptions := mqcloudService.NewSetCertificateAmsChannelsOptions(
				config["QUEUE_MANAGER_ID"],
				config["KEY_STORE_CERTIFICATE_ID"],
				config["SERVICE_INSTANCE_DEPLOYMENT_GUID"],
				[]mqcloudv1.ChannelDetails{*channelDetailsModel},
			)

			setCertificateAmsChannelsOptions.UpdateStrategy = core.StringPtr("replace")

			channelsDetails, response, err := mqcloudService.SetCertificateAmsChannels(setCertificateAmsChannelsOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(channelsDetails, "", "  ")
			fmt.Println(string(b))

			// end-set_certificate_ams_channels

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(channelsDetails).ToNot(BeNil())
		})
		It(`CreateVirtualPrivateEndpointGateway request example`, func() {
			fmt.Println("\nCreateVirtualPrivateEndpointGateway() result:")
			// begin-create_virtual_private_endpoint_gateway

			createVirtualPrivateEndpointGatewayOptions := mqcloudService.NewCreateVirtualPrivateEndpointGatewayOptions(
				config["SERVICE_INSTANCE_CAPACITY_GUID"],
				"vpe_gateway1-to-vpe_gateway2",
				config["TARGET_CRN"],
			)
			createVirtualPrivateEndpointGatewayOptions.SetTrustedProfile(config["TRUSTED_PROFILE"])

			virtualPrivateEndpointGatewayDetails, response, err := mqcloudService.CreateVirtualPrivateEndpointGateway(createVirtualPrivateEndpointGatewayOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(virtualPrivateEndpointGatewayDetails, "", "  ")
			fmt.Println(string(b))
			config["VIRTUAL_PRIVATE_ENDPOINT_GATEWAY_GUID"] = *virtualPrivateEndpointGatewayDetails.ID
			// end-create_virtual_private_endpoint_gateway

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(201))
			Expect(virtualPrivateEndpointGatewayDetails).ToNot(BeNil())
		})
		It(`ListVirtualPrivateEndpointGateways request example`, func() {
			fmt.Println("\nListVirtualPrivateEndpointGateways() result:")
			// begin-list_virtual_private_endpoint_gateways
			listVirtualPrivateEndpointGatewaysOptions := &mqcloudv1.ListVirtualPrivateEndpointGatewaysOptions{
				ServiceInstanceGuid: core.StringPtr(config["SERVICE_INSTANCE_CAPACITY_GUID"]),
				TrustedProfile:      core.StringPtr(config["TRUSTED_PROFILE"]),
				Limit:               core.Int64Ptr(int64(10)),
			}

			pager, err := mqcloudService.NewVirtualPrivateEndpointGatewaysPager(listVirtualPrivateEndpointGatewaysOptions)
			if err != nil {
				panic(err)
			}

			var allResults []mqcloudv1.VirtualPrivateEndpointGatewayDetails
			for pager.HasNext() {
				nextPage, err := pager.GetNext()
				if err != nil {
					panic(err)
				}
				allResults = append(allResults, nextPage...)
			}
			b, _ := json.MarshalIndent(allResults, "", "  ")
			fmt.Println(string(b))
			// end-list_virtual_private_endpoint_gateways
		})
		It(`GetVirtualPrivateEndpointGateway request example`, func() {
			fmt.Println("\nGetVirtualPrivateEndpointGateway() result:")
			// begin-get_virtual_private_endpoint_gateway

			getVirtualPrivateEndpointGatewayOptions := mqcloudService.NewGetVirtualPrivateEndpointGatewayOptions(
				config["SERVICE_INSTANCE_CAPACITY_GUID"],
				config["VIRTUAL_PRIVATE_ENDPOINT_GATEWAY_GUID"],
			)
			getVirtualPrivateEndpointGatewayOptions.SetTrustedProfile(config["TRUSTED_PROFILE"])
			virtualPrivateEndpointGatewayDetails, response, err := mqcloudService.GetVirtualPrivateEndpointGateway(getVirtualPrivateEndpointGatewayOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(virtualPrivateEndpointGatewayDetails, "", "  ")
			fmt.Println(string(b))

			// end-get_virtual_private_endpoint_gateway

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(virtualPrivateEndpointGatewayDetails).ToNot(BeNil())
		})
		It(`DeleteVirtualPrivateEndpointGateway request example`, func() {
			// begin-delete_virtual_private_endpoint_gateway

			deleteVirtualPrivateEndpointGatewayOptions := mqcloudService.NewDeleteVirtualPrivateEndpointGatewayOptions(
				config["SERVICE_INSTANCE_CAPACITY_GUID"],
				config["VIRTUAL_PRIVATE_ENDPOINT_GATEWAY_GUID"],
			)
			deleteVirtualPrivateEndpointGatewayOptions.SetTrustedProfile(config["TRUSTED_PROFILE"])
			response, err := mqcloudService.DeleteVirtualPrivateEndpointGateway(deleteVirtualPrivateEndpointGatewayOptions)
			if err != nil {
				panic(err)
			}
			if response.StatusCode != 204 {
				fmt.Printf("\nUnexpected response status code received from DeleteVirtualPrivateEndpointGateway(): %d\n", response.StatusCode)
			}

			// end-delete_virtual_private_endpoint_gateway

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(204))
		})
		It(`DeleteUser request example`, func() {
			// begin-delete_user

			deleteUserOptions := mqcloudService.NewDeleteUserOptions(
				config["SERVICE_INSTANCE_DEPLOYMENT_GUID"],
				config["USER_ID"],
			)

			response, err := mqcloudService.DeleteUser(deleteUserOptions)
			if err != nil {
				panic(err)
			}
			if response.StatusCode != 204 {
				fmt.Printf("\nUnexpected response status code received from DeleteUser(): %d\n", response.StatusCode)
			}

			// end-delete_user

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(204))
		})
		It(`DeleteApplication request example`, func() {
			// begin-delete_application

			deleteApplicationOptions := mqcloudService.NewDeleteApplicationOptions(
				config["SERVICE_INSTANCE_DEPLOYMENT_GUID"],
				config["APPLICATION_ID"],
			)

			response, err := mqcloudService.DeleteApplication(deleteApplicationOptions)
			if err != nil {
				panic(err)
			}
			if response.StatusCode != 204 {
				fmt.Printf("\nUnexpected response status code received from DeleteApplication(): %d\n", response.StatusCode)
			}

			// end-delete_application

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(204))
		})
		It(`DeleteTrustStoreCertificate request example`, func() {
			// begin-delete_trust_store_certificate

			deleteTrustStoreCertificateOptions := mqcloudService.NewDeleteTrustStoreCertificateOptions(
				config["SERVICE_INSTANCE_DEPLOYMENT_GUID"],
				config["QUEUE_MANAGER_ID"],
				config["TRUST_STORE_CERTIFICATE_ID"],
			)

			response, err := mqcloudService.DeleteTrustStoreCertificate(deleteTrustStoreCertificateOptions)
			if err != nil {
				panic(err)
			}
			if response.StatusCode != 204 {
				fmt.Printf("\nUnexpected response status code received from DeleteTrustStoreCertificate(): %d\n", response.StatusCode)
			}

			// end-delete_trust_store_certificate

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(204))
		})
		It(`DeleteKeyStoreCertificate request example`, func() {
			// begin-delete_key_store_certificate

			deleteKeyStoreCertificateOptions := mqcloudService.NewDeleteKeyStoreCertificateOptions(
				config["SERVICE_INSTANCE_DEPLOYMENT_GUID"],
				config["QUEUE_MANAGER_ID"],
				config["KEY_STORE_CERTIFICATE_ID"],
			)

			response, err := mqcloudService.DeleteKeyStoreCertificate(deleteKeyStoreCertificateOptions)
			if err != nil {
				panic(err)
			}
			if response.StatusCode != 204 {
				fmt.Printf("\nUnexpected response status code received from DeleteKeyStoreCertificate(): %d\n", response.StatusCode)
			}

			// end-delete_key_store_certificate

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(204))
		})
		It(`DeleteQueueManager request example`, func() {
			fmt.Println("\nDeleteQueueManager() result:")
			// begin-delete_queue_manager

			deleteQueueManagerOptions := mqcloudService.NewDeleteQueueManagerOptions(
				config["SERVICE_INSTANCE_DEPLOYMENT_GUID"],
				config["QUEUE_MANAGER_ID"],
			)

			queueManagerTaskStatus, response, err := mqcloudService.DeleteQueueManager(deleteQueueManagerOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(queueManagerTaskStatus, "", "  ")
			fmt.Println(string(b))

			// end-delete_queue_manager

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(202))
			Expect(queueManagerTaskStatus).ToNot(BeNil())
		})

	})
})
