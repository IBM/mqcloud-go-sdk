//go:build integration
// +build integration

/**
 * (C) Copyright IBM Corp. 2023.
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
	const externalConfigFile = "../ibm-credentials.env"

	var (
		err            error
		mqcloudService *mqcloudv1.MqcloudV1
		serviceURL     string

		queue_manager_id          *string
		user_id                   *string
		application_id            *string
		truststore_certificate_id *string
		keystore_certificate_id   *string
		serviceinstance_guid      string
		keystore_filepath         string
		truststore_filepath       string
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

			serviceURL = os.Getenv("IBMCLOUD_MQCLOUD_CONFIG_ENDPOINT")
			if serviceURL == "" {
				Skip("Unable to load service URL configuration property, skipping tests")
			}
			serviceinstance_guid = os.Getenv("SERVICEINSTANCE_GUID")
			keystore_filepath = os.Getenv("KEYSTORE_FILE_PATH")
			truststore_filepath = os.Getenv("TRUSTSTORE_FILE_PATH")
			fmt.Fprintf(GinkgoWriter, "Service URL: %v\n", serviceURL)
			shouldSkipTest = func() {}
		})
	})

	Describe(`Client initialization`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It("Successfully construct the service client instance", func() {
			apikey := os.Getenv("IC_API_KEY")
			authenticator := &core.IamAuthenticator{
				ApiKey: apikey,
				URL:    "https://iam.test.cloud.ibm.com" + "/identity/token",
			}
			if err != nil {
				panic(err)
			}
			mqcloudServiceOptions := &mqcloudv1.MqcloudV1Options{
				URL:           os.Getenv("IBMCLOUD_MQCLOUD_CONFIG_ENDPOINT"),
				Authenticator: authenticator,
				ServiceName:   mqcloudv1.DefaultServiceName,
			}

			mqcloudService, err = mqcloudv1.NewMqcloudV1UsingExternalConfig(mqcloudServiceOptions)
			Expect(err).To(BeNil())
			Expect(mqcloudService).ToNot(BeNil())
			Expect(mqcloudService.Service.Options.URL).To(Equal(serviceURL))

			core.SetLogger(core.NewLogger(core.LevelDebug, log.New(GinkgoWriter, "", log.LstdFlags), log.New(GinkgoWriter, "", log.LstdFlags)))
			mqcloudService.EnableRetries(4, 30*time.Second)
		})
	})

	// Describe(`GetUsageDetails - Get the usage details`, func() {
	// 	BeforeEach(func() {
	// 		shouldSkipTest()
	// 	})
	// 	It(`GetUsageDetails(getUsageDetailsOptions *GetUsageDetailsOptions)`, func() {
	// 		fmt.Println("GetusageDetails.....")
	// 		getUsageDetailsOptions := &mqcloudv1.GetUsageDetailsOptions{
	// 			ServiceInstanceGuid: core.StringPtr("7256e811-f88e-46b5-bbc2-05b8a8adcf3e"),
	// 		}

	// 		usage, response, err := mqcloudService.GetUsageDetails(getUsageDetailsOptions)
	// 		Expect(err).To(BeNil())
	// 		Expect(response.StatusCode).To(Equal(200))
	// 		Expect(usage).To(BeNil())
	// 	})
	// })

	// Describe(`GetOptions - Return configuration options (eg, available deployment locations, queue manager sizes)`, func() {
	// 	BeforeEach(func() {
	// 		shouldSkipTest()
	// 	})
	// 	It(`GetOptions(getOptionsOptions *GetOptionsOptions)`, func() {
	// 		getOptionsOptions := &mqcloudv1.GetOptionsOptions{
	// 			ServiceInstanceGuid: core.StringPtr(serviceinstance_guid),
	// 		}

	// 		configurationOptions, response, err := mqcloudService.GetOptions(getOptionsOptions)
	// 		Expect(err).To(BeNil())
	// 		Expect(response.StatusCode).To(Equal(200))
	// 		Expect(configurationOptions).ToNot(BeNil())
	// 	})
	// })

	Describe(`CreateQueueManager - Create a new queue manager`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`CreateQueueManager(createQueueManagerOptions *CreateQueueManagerOptions)`, func() {
			createQueueManagerOptions := &mqcloudv1.CreateQueueManagerOptions{
				ServiceInstanceGuid: core.StringPtr(serviceinstance_guid),
				Name:                core.StringPtr("int_test" + RandString(6)),
				Location:            core.StringPtr("reserved-eu-de-cluster-f884"),
				Size:                core.StringPtr("lite"),
				DisplayName:         core.StringPtr("A test queue manager"),
				Version:             core.StringPtr("9.3.3_3"),
			}

			queueManagerTaskStatus, response, err := mqcloudService.CreateQueueManager(createQueueManagerOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(202))
			Expect(queueManagerTaskStatus).ToNot(BeNil())
			queue_manager_id = queueManagerTaskStatus.QueueManagerID
		})
	})
	Describe(`ListQueueManagers - Get list of queue managers`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`ListQueueManagers(listQueueManagersOptions *ListQueueManagersOptions) with pagination`, func() {
			listQueueManagersOptions := &mqcloudv1.ListQueueManagersOptions{
				ServiceInstanceGuid: core.StringPtr(serviceinstance_guid),
				AcceptLanguage:      core.StringPtr("testString"),
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
				ServiceInstanceGuid: core.StringPtr(serviceinstance_guid),
				AcceptLanguage:      core.StringPtr("testString"),
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
		})
		It(`GetQueueManager(getQueueManagerOptions *GetQueueManagerOptions)`, func() {
			getQueueManagerOptions := &mqcloudv1.GetQueueManagerOptions{
				ServiceInstanceGuid: core.StringPtr(serviceinstance_guid),
				QueueManagerID:      queue_manager_id,
				AcceptLanguage:      core.StringPtr("testString"),
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
		})
		It(`SetQueueManagerVersion(setQueueManagerVersionOptions *SetQueueManagerVersionOptions)`, func() {
			WaitForQmStatusUpdate(queue_manager_id, mqcloudService, serviceinstance_guid)
			fmt.Fprintf(GinkgoWriter,
				"--------- Queue Manager is now in the running state ---------",
			)
			setQueueManagerVersionOptions := &mqcloudv1.SetQueueManagerVersionOptions{
				ServiceInstanceGuid: core.StringPtr(serviceinstance_guid),
				QueueManagerID:      queue_manager_id,
				Version:             core.StringPtr("9.3.4_1"),
				AcceptLanguage:      core.StringPtr("testString"),
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
				ServiceInstanceGuid: core.StringPtr(serviceinstance_guid),
				QueueManagerID:      queue_manager_id,
				AcceptLanguage:      core.StringPtr("testString"),
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
				ServiceInstanceGuid: core.StringPtr(serviceinstance_guid),
				QueueManagerID:      queue_manager_id,
				AcceptLanguage:      core.StringPtr("testString"),
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
				ServiceInstanceGuid: core.StringPtr(serviceinstance_guid),
				QueueManagerID:      queue_manager_id,
				AcceptLanguage:      core.StringPtr("testString"),
			}

			queueManagerStatus, response, err := mqcloudService.GetQueueManagerStatus(getQueueManagerStatusOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(queueManagerStatus).ToNot(BeNil())
		})
	})

	Describe(`CreateUser - Add a user to an instance`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`CreateUser(createUserOptions *CreateUserOptions)`, func() {
			createUserOptions := &mqcloudv1.CreateUserOptions{
				ServiceInstanceGuid: core.StringPtr(serviceinstance_guid),
				Email:               core.StringPtr("user" + RandString(6) + "@ibm.com"),
				Name:                core.StringPtr("user" + RandString(6)),
				AcceptLanguage:      core.StringPtr("testString"),
			}

			userDetails, response, err := mqcloudService.CreateUser(createUserOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(201))
			Expect(userDetails).ToNot(BeNil())
			user_id = userDetails.ID
		})
	})

	Describe(`GetUser - Get a user for an instance`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`GetUser(getUserOptions *GetUserOptions)`, func() {
			getUserOptions := &mqcloudv1.GetUserOptions{
				ServiceInstanceGuid: core.StringPtr(serviceinstance_guid),
				UserID:              user_id,
				AcceptLanguage:      core.StringPtr("testString"),
			}

			userDetails, response, err := mqcloudService.GetUser(getUserOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(userDetails).ToNot(BeNil())
		})
	})
	Describe(`ListUsers - Get a list of users for an instance`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`ListUsers(listUsersOptions *ListUsersOptions) with pagination`, func() {
			listUsersOptions := &mqcloudv1.ListUsersOptions{
				ServiceInstanceGuid: core.StringPtr(serviceinstance_guid),
				AcceptLanguage:      core.StringPtr("testString"),
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
				ServiceInstanceGuid: core.StringPtr(serviceinstance_guid),
				AcceptLanguage:      core.StringPtr("testString"),
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

	Describe(`CreateApplication - Add an application to an instance`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`CreateApplication(createApplicationOptions *CreateApplicationOptions)`, func() {
			createApplicationOptions := &mqcloudv1.CreateApplicationOptions{
				ServiceInstanceGuid: core.StringPtr(serviceinstance_guid),
				Name:                core.StringPtr("app" + RandString(1)),
				AcceptLanguage:      core.StringPtr("testString"),
			}

			applicationCreated, response, err := mqcloudService.CreateApplication(createApplicationOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(201))
			Expect(applicationCreated).ToNot(BeNil())
			application_id = applicationCreated.ID
		})
	})

	Describe(`GetApplication - Get an application for an instance`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`GetApplication(getApplicationOptions *GetApplicationOptions)`, func() {
			getApplicationOptions := &mqcloudv1.GetApplicationOptions{
				ServiceInstanceGuid: core.StringPtr(serviceinstance_guid),
				ApplicationID:       application_id,
				AcceptLanguage:      core.StringPtr("testString"),
			}

			applicationDetails, response, err := mqcloudService.GetApplication(getApplicationOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(applicationDetails).ToNot(BeNil())
		})
	})

	Describe(`ListApplications - Get a list of applications for an instance`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`ListApplications(listApplicationsOptions *ListApplicationsOptions) with pagination`, func() {
			listApplicationsOptions := &mqcloudv1.ListApplicationsOptions{
				ServiceInstanceGuid: core.StringPtr(serviceinstance_guid),
				AcceptLanguage:      core.StringPtr("testString"),
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
		It(`ListApplications(listApplicationsOptions *ListApplicationsOptions) using ApplicationsPager`, func() {
			listApplicationsOptions := &mqcloudv1.ListApplicationsOptions{
				ServiceInstanceGuid: core.StringPtr(serviceinstance_guid),
				AcceptLanguage:      core.StringPtr("testString"),
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

	Describe(`CreateApplicationApikey - Create a new apikey for an application`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`CreateApplicationApikey(createApplicationApikeyOptions *CreateApplicationApikeyOptions)`, func() {
			createApplicationApikeyOptions := &mqcloudv1.CreateApplicationApikeyOptions{
				ServiceInstanceGuid: core.StringPtr(serviceinstance_guid),
				ApplicationID:       application_id,
				Name:                core.StringPtr("test-api-key"),
				AcceptLanguage:      core.StringPtr("testString"),
			}

			applicationApiKeyCreated, response, err := mqcloudService.CreateApplicationApikey(createApplicationApikeyOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(201))
			Expect(applicationApiKeyCreated).ToNot(BeNil())
		})
	})

	Describe(`CreateTrustStorePemCertificate - Upload a certificate`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`CreateTrustStorePemCertificate(createTrustStorePemCertificateOptions *CreateTrustStorePemCertificateOptions)`, func() {
			WaitForQmStatusUpdate(queue_manager_id, mqcloudService, serviceinstance_guid)
			filePath := truststore_filepath // Replace with your file path
			file, err := os.Open(filePath)
			if err != nil {
				fmt.Println(GinkgoWriter, fmt.Errorf("Error opening file: %s", err))
				return
			}
			defer file.Close()
			createTrustStorePemCertificateOptions := &mqcloudv1.CreateTrustStorePemCertificateOptions{
				ServiceInstanceGuid: core.StringPtr(serviceinstance_guid),
				QueueManagerID:      queue_manager_id,
				Label:               core.StringPtr("testString101"),
				CertificateFile:     file,
				AcceptLanguage:      core.StringPtr("testString"),
			}

			trustStoreCertificateDetails, response, err := mqcloudService.CreateTrustStorePemCertificate(createTrustStorePemCertificateOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(201))
			Expect(trustStoreCertificateDetails).ToNot(BeNil())
			truststore_certificate_id = trustStoreCertificateDetails.ID
		})
	})

	Describe(`ListTrustStoreCertificates - List certificates`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`ListTrustStoreCertificates(listTrustStoreCertificatesOptions *ListTrustStoreCertificatesOptions)`, func() {
			listTrustStoreCertificatesOptions := &mqcloudv1.ListTrustStoreCertificatesOptions{
				ServiceInstanceGuid: core.StringPtr(serviceinstance_guid),
				QueueManagerID:      queue_manager_id,
				AcceptLanguage:      core.StringPtr("testString"),
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
				ServiceInstanceGuid: core.StringPtr(serviceinstance_guid),
				QueueManagerID:      queue_manager_id,
				CertificateID:       truststore_certificate_id,
				AcceptLanguage:      core.StringPtr("testString"),
			}

			trustStoreCertificateDetails, response, err := mqcloudService.GetTrustStoreCertificate(getTrustStoreCertificateOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(trustStoreCertificateDetails).ToNot(BeNil())
		})
	})

	Describe(`DeleteTrustStoreCertificate - Delete a trust store certificate`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		// time.Sleep(1 * time.Minute)
		It(`DeleteTrustStoreCertificate(deleteTrustStoreCertificateOptions *DeleteTrustStoreCertificateOptions)`, func() {
			deleteTrustStoreCertificateOptions := &mqcloudv1.DeleteTrustStoreCertificateOptions{
				ServiceInstanceGuid: core.StringPtr(serviceinstance_guid),
				QueueManagerID:      queue_manager_id,
				CertificateID:       truststore_certificate_id,
				AcceptLanguage:      core.StringPtr("testString"),
			}

			response, err := mqcloudService.DeleteTrustStoreCertificate(deleteTrustStoreCertificateOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(204))
		})
	})

	// Describe(`DownloadTrustStoreCertificate - Download a queue manager's certificate from its trust store`, func() {
	// 	BeforeEach(func() {
	// 		shouldSkipTest()
	// 	})
	// 	It(`DownloadTrustStoreCertificate(downloadTrustStoreCertificateOptions *DownloadTrustStoreCertificateOptions)`, func() {
	// 		downloadTrustStoreCertificateOptions := &mqcloudv1.DownloadTrustStoreCertificateOptions{
	// 			ServiceInstanceGuid: core.StringPtr(serviceinstance_guid),
	// 			QueueManagerID:      queue_manager_id,
	// 			CertificateID:       truststore_certificate_id,
	// 			AcceptLanguage:      core.StringPtr("testString"),
	// 		}

	// 		result, response, err := mqcloudService.DownloadTrustStoreCertificate(downloadTrustStoreCertificateOptions)
	// 		Expect(err).To(BeNil())
	// 		Expect(response.StatusCode).To(Equal(200))
	// 		Expect(result).ToNot(BeNil())
	// 	})
	// })

	Describe(`CreateKeyStorePemCertificate - Upload a certificate`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`CreateKeyStorePemCertificate(createKeyStorePemCertificateOptions *CreateKeyStorePemCertificateOptions)`, func() {
			filePath := keystore_filepath // Replace with your file path
			file, err := os.Open(filePath)
			if err != nil {
				fmt.Println(GinkgoWriter, fmt.Errorf("Error opening file: %s", err))
				return
			}
			defer file.Close()
			createKeyStorePemCertificateOptions := &mqcloudv1.CreateKeyStorePemCertificateOptions{
				ServiceInstanceGuid: core.StringPtr(serviceinstance_guid),
				QueueManagerID:      queue_manager_id,
				Label:               core.StringPtr("testString"),
				CertificateFile:     file,
				AcceptLanguage:      core.StringPtr("testString"),
			}

			keyStoreCertificateDetails, response, err := mqcloudService.CreateKeyStorePemCertificate(createKeyStorePemCertificateOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(201))
			Expect(keyStoreCertificateDetails).ToNot(BeNil())
			keystore_certificate_id = keyStoreCertificateDetails.ID
		})
	})

	Describe(`ListKeyStoreCertificates - List certificates`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`ListKeyStoreCertificates(listKeyStoreCertificatesOptions *ListKeyStoreCertificatesOptions)`, func() {
			listKeyStoreCertificatesOptions := &mqcloudv1.ListKeyStoreCertificatesOptions{
				ServiceInstanceGuid: core.StringPtr(serviceinstance_guid),
				QueueManagerID:      queue_manager_id,
				AcceptLanguage:      core.StringPtr("testString"),
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
				ServiceInstanceGuid: core.StringPtr(serviceinstance_guid),
				QueueManagerID:      queue_manager_id,
				CertificateID:       keystore_certificate_id,
				AcceptLanguage:      core.StringPtr("testString"),
			}

			keyStoreCertificateDetails, response, err := mqcloudService.GetKeyStoreCertificate(getKeyStoreCertificateOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(keyStoreCertificateDetails).ToNot(BeNil())
		})
	})

	// Describe(`DownloadKeyStoreCertificate - Download a queue manager's certificate from its key store`, func() {
	// 	BeforeEach(func() {
	// 		shouldSkipTest()
	// 	})
	// 	It(`DownloadKeyStoreCertificate(downloadKeyStoreCertificateOptions *DownloadKeyStoreCertificateOptions)`, func() {
	// 		downloadKeyStoreCertificateOptions := &mqcloudv1.DownloadKeyStoreCertificateOptions{
	// 			ServiceInstanceGuid: core.StringPtr(serviceinstance_guid),
	// 			QueueManagerID:      queue_manager_id,
	// 			CertificateID:       keystore_certificate_id,
	// 			AcceptLanguage:      core.StringPtr("testString"),
	// 		}

	// 		result, response, err := mqcloudService.DownloadKeyStoreCertificate(downloadKeyStoreCertificateOptions)
	// 		Expect(err).To(BeNil())
	// 		Expect(response.StatusCode).To(Equal(200))
	// 		Expect(result).ToNot(BeNil())
	// 	})
	// })

	// Describe(`GetCertificateAmsChannels - Get the AMS channels that are configured with this key store certificate`, func() {
	// 	BeforeEach(func() {
	// 		shouldSkipTest()
	// 	})
	// 	It(`GetCertificateAmsChannels(getCertificateAmsChannelsOptions *GetCertificateAmsChannelsOptions)`, func() {
	// 		getCertificateAmsChannelsOptions := &mqcloudv1.GetCertificateAmsChannelsOptions{
	// 			QueueManagerID:      queue_manager_id,
	// 			CertificateID:       keystore_certificate_id,
	// 			ServiceInstanceGuid: core.StringPtr(serviceinstance_guid),
	// 			AcceptLanguage:      core.StringPtr("testString"),
	// 		}

	// 		channelsDetails, response, err := mqcloudService.GetCertificateAmsChannels(getCertificateAmsChannelsOptions)
	// 		Expect(err).To(BeNil())
	// 		Expect(response.StatusCode).To(Equal(200))
	// 		Expect(channelsDetails).ToNot(BeNil())
	// 	})
	// })

	// Describe(`SetCertificateAmsChannels - Update the AMS channels that are configured with this key store certificate`, func() {
	// 	BeforeEach(func() {
	// 		shouldSkipTest()
	// 	})
	// 	It(`SetCertificateAmsChannels(setCertificateAmsChannelsOptions *SetCertificateAmsChannelsOptions)`, func() {
	// 		channelDetailsModel := &mqcloudv1.ChannelDetails{
	// 			Name: core.StringPtr("testString"),
	// 		}

	// 		setCertificateAmsChannelsOptions := &mqcloudv1.SetCertificateAmsChannelsOptions{
	// 			QueueManagerID:      queue_manager_id,
	// 			CertificateID:       keystore_certificate_id,
	// 			ServiceInstanceGuid: core.StringPtr(serviceinstance_guid),
	// 			Channels:            []mqcloudv1.ChannelDetails{*channelDetailsModel},
	// 			UpdateStrategy:      core.StringPtr("replace"),
	// 			AcceptLanguage:      core.StringPtr("testString"),
	// 		}

	// 		channelsDetails, response, err := mqcloudService.SetCertificateAmsChannels(setCertificateAmsChannelsOptions)
	// 		Expect(err).To(BeNil())
	// 		Expect(response.StatusCode).To(Equal(200))
	// 		Expect(channelsDetails).ToNot(BeNil())
	// 	})
	// })

	Describe(`DeleteKeyStoreCertificate - Delete a queue manager's key store certificate`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`DeleteKeyStoreCertificate(deleteKeyStoreCertificateOptions *DeleteKeyStoreCertificateOptions)`, func() {
			deleteKeyStoreCertificateOptions := &mqcloudv1.DeleteKeyStoreCertificateOptions{
				ServiceInstanceGuid: core.StringPtr(serviceinstance_guid),
				QueueManagerID:      queue_manager_id,
				CertificateID:       keystore_certificate_id,
				AcceptLanguage:      core.StringPtr("testString"),
			}

			response, err := mqcloudService.DeleteKeyStoreCertificate(deleteKeyStoreCertificateOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(204))
		})
	})

	Describe(`DeleteUser - Delete a user for an instance`, func() {
		BeforeEach(func() {
			shouldSkipTest()
		})
		It(`DeleteUser(deleteUserOptions *DeleteUserOptions)`, func() {
			deleteUserOptions := &mqcloudv1.DeleteUserOptions{
				ServiceInstanceGuid: core.StringPtr(serviceinstance_guid),
				UserID:              user_id,
				AcceptLanguage:      core.StringPtr("testString"),
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
				ServiceInstanceGuid: core.StringPtr(serviceinstance_guid),
				ApplicationID:       application_id,
				AcceptLanguage:      core.StringPtr("testString"),
			}

			response, err := mqcloudService.DeleteApplication(deleteApplicationOptions)
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
				ServiceInstanceGuid: core.StringPtr(serviceinstance_guid),
				QueueManagerID:      queue_manager_id,
				AcceptLanguage:      core.StringPtr("testString"),
			}

			queueManagerTaskStatus, response, err := mqcloudService.DeleteQueueManager(deleteQueueManagerOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(202))
			Expect(queueManagerTaskStatus).ToNot(BeNil())
		})
	})
})

//
// Utility functions are declared in the unit test file
//
