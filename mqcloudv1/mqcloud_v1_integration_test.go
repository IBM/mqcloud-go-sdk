//go:build integration

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
	"fmt"
	"log"
	"os"
	"strconv"
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

	Describe(`External configuration`, func() {
		It("Successfully load the configuration", func() {
			_, err = os.Stat(externalConfigFile)
			if err != nil {
				Skip("External configuration file not found, skipping tests: " + err.Error())
			}

			_ = os.Setenv("IBM_CREDENTIALS_FILE", externalConfigFile)
			config, err = core.GetServiceProperties(mqcloudv1.DefaultServiceName)
			if err != nil {
				Skip("Error loading service properties, skipping tests: " + err.Error())
			}
			serviceURL = config["URL"]
			if serviceURL == "" {
				Skip("Unable to load service URL configuration property, skipping tests")
			}

			_, _ = fmt.Fprintf(GinkgoWriter, "Service URL: %v\n", serviceURL)
		})
	})

	Describe(`Client initialization`, func() {
		BeforeEach(func() {
			SkipUnlessConfigured(config, "URL")
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
			SkipUnlessConfigured(config, "SERVICE_INSTANCE_DEPLOYMENT_GUID")
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
			SkipUnlessConfigured(config, "SERVICE_INSTANCE_DEPLOYMENT_GUID")
		})
		It(`GetOptions(getOptionsOptions *GetOptionsOptions)`, func() {
			getOptionsOptions := &mqcloudv1.GetOptionsOptions{
				ServiceInstanceGuid: core.StringPtr(config["SERVICE_INSTANCE_DEPLOYMENT_GUID"]),
			}

			configurationOptions, response, err := mqcloudService.GetOptions(getOptionsOptions)
			Expect(err).To(BeNil())
			config["LOCATION"] = configurationOptions.Locations[0]
			config["VERSION"], config["VERSION_UPGRADE"] = ChooseQueueManagerVersions(
				configurationOptions.Versions,
				config["VERSION"],
				config["VERSION_UPGRADE"],
			)
			Expect(response.StatusCode).To(Equal(200))
			Expect(configurationOptions).ToNot(BeNil())
		})
	})

	Describe(`CreateQueueManager - Create a new queue manager`, func() {
		BeforeEach(func() {
			SkipUnlessConfigured(config, "SERVICE_INSTANCE_DEPLOYMENT_GUID", "LOCATION", "VERSION")
		})
		It(`CreateQueueManager(createQueueManagerOptions *CreateQueueManagerOptions)`, func() {
			createQueueManagerOptions := &mqcloudv1.CreateQueueManagerOptions{
				ServiceInstanceGuid: core.StringPtr(config["SERVICE_INSTANCE_DEPLOYMENT_GUID"]),
				Name:                core.StringPtr("testqm22"),
				Location:            core.StringPtr(config["LOCATION"]),
				Size:                core.StringPtr("small"),
				DisplayName:         core.StringPtr("A test queue manager"),
				Version:             core.StringPtr(config["VERSION"]),
			}

			queueManagerTaskStatus, response, err := mqcloudService.CreateQueueManager(createQueueManagerOptions)
			Expect(err).To(BeNil())
			config["QUEUE_MANAGER_ID"] = *queueManagerTaskStatus.QueueManagerID
			Expect(response.StatusCode).To(Equal(202))
			Expect(queueManagerTaskStatus).ToNot(BeNil())
			// Wait for queue manager to be running
			err = WaitForQueueManagerRunning(mqcloudService, config["SERVICE_INSTANCE_DEPLOYMENT_GUID"], config["QUEUE_MANAGER_ID"])
			Expect(err).To(BeNil())
		})
	})

	Describe(`ListQueueManagers - Get list of queue managers`, func() {
		BeforeEach(func() {
			SkipUnlessConfigured(config, "SERVICE_INSTANCE_DEPLOYMENT_GUID")
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
			_, _ = fmt.Fprintf(GinkgoWriter, "Retrieved a total of %d item(s) with pagination.\n", len(allResults))
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
			_, _ = fmt.Fprintf(GinkgoWriter, "ListQueueManagers() returned a total of %d item(s) using QueueManagersPager.\n", len(allResults))
		})
	})

	Describe(`GetQueueManager - Get details of a queue manager`, func() {
		BeforeEach(func() {
			SkipTestIfQmIsNotRunning(config, mqcloudService)
		})
		It(`GetQueueManager(getQueueManagerOptions *GetQueueManagerOptions)`, func() {
			_, _ = fmt.Fprintln(GinkgoWriter, "🚀 Starting test: GetQueueManager")

			getQueueManagerOptions := &mqcloudv1.GetQueueManagerOptions{
				ServiceInstanceGuid: core.StringPtr(config["SERVICE_INSTANCE_DEPLOYMENT_GUID"]),
				QueueManagerID:      core.StringPtr(config["QUEUE_MANAGER_ID"]),
			}

			queueManagerDetails, response, err := mqcloudService.GetQueueManager(getQueueManagerOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(queueManagerDetails).ToNot(BeNil())

			// Extract and log details using nil checks
			id := "-"
			if queueManagerDetails.ID != nil {
				id = *queueManagerDetails.ID
			}
			name := "-"
			if queueManagerDetails.Name != nil {
				name = *queueManagerDetails.Name
			}
			displayName := "-"
			if queueManagerDetails.DisplayName != nil {
				displayName = *queueManagerDetails.DisplayName
			}
			location := "-"
			if queueManagerDetails.Location != nil {
				location = *queueManagerDetails.Location
			}
			size := "-"
			if queueManagerDetails.Size != nil {
				size = *queueManagerDetails.Size
			}
			version := "-"
			if queueManagerDetails.Version != nil {
				version = *queueManagerDetails.Version
			}
			webConsole := "-"
			if queueManagerDetails.WebConsoleURL != nil {
				webConsole = *queueManagerDetails.WebConsoleURL
			}
			created := "-"
			if queueManagerDetails.DateCreated != nil {
				created = queueManagerDetails.DateCreated.String()
			}

			fmt.Printf("✅ Queue Manager [%s] (Display Name: %s)\n", name, displayName)
			fmt.Printf("🆔 ID: %s\n🌍 Location: %s\n📦 Size: %s\n📈 Version: %s\n🌐 Web Console: %s\n🕒 Created: %s\n",
				id, location, size, version, webConsole, created)

			time.Sleep(10 * time.Second)
		})

	})

	Describe(`SetQueueManagerVersion - Upgrade a queue manager`, func() {
		BeforeEach(func() {
			SkipTestIfQmIsNotRunning(config, mqcloudService)
			SkipUnlessConfigured(config, "VERSION_UPGRADE")
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

			err = WaitForQueueManagerVersion(mqcloudService, config["SERVICE_INSTANCE_DEPLOYMENT_GUID"], config["QUEUE_MANAGER_ID"], config["VERSION_UPGRADE"])
			Expect(err).To(BeNil())
		})
	})

	Describe(`GetQueueManagerAvailableUpgradeVersions - Get the list of available versions that this queue manager can be upgraded to`, func() {
		BeforeEach(func() {
			SkipUnlessConfigured(config, "SERVICE_INSTANCE_DEPLOYMENT_GUID", "QUEUE_MANAGER_ID")
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
			SkipUnlessConfigured(config, "SERVICE_INSTANCE_DEPLOYMENT_GUID", "QUEUE_MANAGER_ID")
		})

		It(`GetQueueManagerConnectionInfo(getQueueManagerConnectionInfoOptions *GetQueueManagerConnectionInfoOptions)`, func() {
			_, _ = fmt.Fprintln(GinkgoWriter, "🔌 Starting test: GetQueueManagerConnectionInfo")

			getQueueManagerConnectionInfoOptions := &mqcloudv1.GetQueueManagerConnectionInfoOptions{
				ServiceInstanceGuid: core.StringPtr(config["SERVICE_INSTANCE_DEPLOYMENT_GUID"]),
				QueueManagerID:      core.StringPtr(config["QUEUE_MANAGER_ID"]),
			}

			connectionInfo, response, err := mqcloudService.GetQueueManagerConnectionInfo(getQueueManagerConnectionInfoOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(connectionInfo).ToNot(BeNil())

			fmt.Println("✅ Queue Manager Connection Info retrieved:")
			for i, ch := range connectionInfo.Channel {
				channelName := "-"
				channelType := "-"

				if ch.Name != nil {
					channelName = *ch.Name
				}
				if ch.Type != nil {
					channelType = *ch.Type
				}

				fmt.Printf("🔌 Channel [%d]:\n", i+1)
				fmt.Printf("   📛 Name : %s\n", channelName)
				fmt.Printf("   🧭 Type : %s\n", channelType)
			}
		})
	})

	Describe(`GetQueueManagerStatus - Get the status of the queue manager`, func() {
		BeforeEach(func() {
			SkipUnlessConfigured(config, "SERVICE_INSTANCE_DEPLOYMENT_GUID", "QUEUE_MANAGER_ID")
		})

		It(`GetQueueManagerStatus(getQueueManagerStatusOptions *GetQueueManagerStatusOptions)`, func() {
			_, _ = fmt.Fprintln(GinkgoWriter, "📡 Starting test: GetQueueManagerStatus")

			getQueueManagerStatusOptions := &mqcloudv1.GetQueueManagerStatusOptions{
				ServiceInstanceGuid: core.StringPtr(config["SERVICE_INSTANCE_DEPLOYMENT_GUID"]),
				QueueManagerID:      core.StringPtr(config["QUEUE_MANAGER_ID"]),
			}

			queueManagerStatus, response, err := mqcloudService.GetQueueManagerStatus(getQueueManagerStatusOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(queueManagerStatus).ToNot(BeNil())

			fmt.Println("✅ Queue Manager Status retrieved:")

			if queueManagerStatus.Status != nil {
				fmt.Printf("   🚦 Status: %s\n", *queueManagerStatus.Status)
			} else {
				fmt.Println("   🚦 Status: <nil>")
			}
		})
	})

	Describe(`ListUsers - Get a list of users for an instance`, func() {
		BeforeEach(func() {
			SkipUnlessConfigured(config, "SERVICE_INSTANCE_DEPLOYMENT_GUID")
		})

		It(`ListUsers(listUsersOptions *ListUsersOptions) with pagination`, func() {
			listUsersOptions := &mqcloudv1.ListUsersOptions{
				ServiceInstanceGuid: core.StringPtr(config["SERVICE_INSTANCE_DEPLOYMENT_GUID"]),
				Limit:               core.Int64Ptr(1),
			}

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

			fmt.Printf("📋 Retrieved a total of %d user(s) with pagination:\n", len(allResults))
			for i, user := range allResults {
				fmt.Printf("👤 User [%d]:\n", i+1)
				if user.ID != nil {
					fmt.Printf("   🆔 ID            : %s\n", *user.ID)
				}
				if user.Name != nil {
					fmt.Printf("   📛 Name          : %s\n", *user.Name)
				}
				if user.Email != nil {
					fmt.Printf("   📧 Email         : %s\n", *user.Email)
				}
				if user.IamServiceID != nil {
					fmt.Printf("   🔐 IAM Service ID: %s\n", *user.IamServiceID)
				}
				if user.Href != nil {
					fmt.Printf("   🔗 Href          : %s\n", *user.Href)
				}
				if user.Roles != nil {
					fmt.Printf("   Roles          : %v\n", user.Roles)
				}
				if user.IamManaged != nil {
					fmt.Printf("   IamManaged          : %v\n", user.IamManaged)
				}
			}
		})

		It(`ListUsers(listUsersOptions *ListUsersOptions) using UsersPager`, func() {
			listUsersOptions := &mqcloudv1.ListUsersOptions{
				ServiceInstanceGuid: core.StringPtr(config["SERVICE_INSTANCE_DEPLOYMENT_GUID"]),
				Limit:               core.Int64Ptr(10),
			}

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

			fmt.Printf("📦 Retrieved %d user(s) using UsersPager:\n", len(allResults))
			for i, user := range allResults {
				fmt.Printf("👤 User [%d]:\n", i+1)
				if user.ID != nil {
					fmt.Printf("   🆔 ID            : %s\n", *user.ID)
				}
				if user.Name != nil {
					fmt.Printf("   📛 Name          : %s\n", *user.Name)
				}
				if user.Email != nil {
					fmt.Printf("   📧 Email         : %s\n", *user.Email)
				}
				if user.IamServiceID != nil {
					fmt.Printf("   🔐 IAM Service ID: %s\n", *user.IamServiceID)
				}
				if user.Href != nil {
					fmt.Printf("   🔗 Href          : %s\n", *user.Href)
				}
				if user.Roles != nil {
					fmt.Printf("   Roles          : %v\n", user.Roles)
				}
				if user.IamManaged != nil {
					fmt.Printf("   IamManaged          : %v\n", user.IamManaged)
				}
			}
		})
	})

	Describe(`CreateUser - Add a user to an instance`, func() {
		BeforeEach(func() {
			SkipUnlessConfigured(config, "SERVICE_INSTANCE_DEPLOYMENT_GUID")
		})
		It(`CreateUser(createUserOptions *CreateUserOptions)`, func() {
			createUserOptions := &mqcloudv1.CreateUserOptions{
				ServiceInstanceGuid: core.StringPtr(config["SERVICE_INSTANCE_DEPLOYMENT_GUID"]),
				Name:                core.StringPtr("user-test22"),
				Email:               core.StringPtr("testemail@ibm.com"),
			}

			userDetails, response, err := mqcloudService.CreateUser(createUserOptions)
			Expect(err).To(BeNil())
			config["USER_ID"] = *userDetails.ID
			Expect(response.StatusCode).To(Equal(201))
			Expect(userDetails).ToNot(BeNil())
			fmt.Println("✅ User has been created successfully")
			// Add a pause to ensure QM is fully created
			time.Sleep(30 * time.Second) // Adjust duration as needed
		})
	})

	Describe(`GetUser - Get a user for an instance`, func() {
		BeforeEach(func() {
			SkipUnlessConfigured(config, "SERVICE_INSTANCE_DEPLOYMENT_GUID", "USER_ID")
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

	Describe(`SetUserName - Update the name for a user`, func() {
		BeforeEach(func() {
			SkipUnlessConfigured(config, "SERVICE_INSTANCE_DEPLOYMENT_GUID", "USER_ID")
		})
		It(`SetUserName(setUserNameOptions *SetUserNameOptions)`, func() {
			setUserNameOptions := &mqcloudv1.SetUserNameOptions{
				ServiceInstanceGuid: core.StringPtr(config["SERVICE_INSTANCE_DEPLOYMENT_GUID"]),
				UserID:              core.StringPtr(config["USER_ID"]),
				Name:                core.StringPtr("t2"),
			}

			userDetails, response, err := mqcloudService.SetUserName(setUserNameOptions)
			Expect(err).To(BeNil())
			Expect(response).ToNot(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(userDetails).ToNot(BeNil())

			fmt.Printf("✅ User shortname has been successfully updated to '%s'\n", *setUserNameOptions.Name)
		})
	})

	Describe(`ListApplications - Get a list of applications for an instance`, func() {
		BeforeEach(func() {
			SkipUnlessConfigured(config, "SERVICE_INSTANCE_DEPLOYMENT_GUID")
		})
		It(`ListApplications(listApplicationsOptions *ListApplicationsOptions) with pagination`, func() {
			listApplicationsOptions := &mqcloudv1.ListApplicationsOptions{
				ServiceInstanceGuid: core.StringPtr(config["SERVICE_INSTANCE_DEPLOYMENT_GUID"]),
				Limit:               core.Int64Ptr(1),
			}

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

			fmt.Printf("📦 Retrieved a total of %d application(s) with pagination:\n", len(allResults))
			for i, app := range allResults {
				fmt.Printf("📲 Application [%d]:\n", i+1)
				if app.ID != nil {
					fmt.Printf("   🆔 ID   : %s\n", *app.ID)
				}
				if app.Name != nil {
					fmt.Printf("   📛 Name : %s\n", *app.Name)
				}
				if app.Href != nil {
					fmt.Printf("   🔗 Href : %s\n", *app.Href)
				}
			}
		})
		It(`ListApplications(listApplicationsOptions *ListApplicationsOptions) using ApplicationsPager`, func() {
			listApplicationsOptions := &mqcloudv1.ListApplicationsOptions{
				ServiceInstanceGuid: core.StringPtr(config["SERVICE_INSTANCE_DEPLOYMENT_GUID"]),
				Limit:               core.Int64Ptr(10),
			}

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

			fmt.Printf("📦 Retrieved a total of %d application(s) using ApplicationsPager:\n", len(allResults))
			for i, app := range allResults {
				fmt.Printf("📲 Application [%d]:\n", i+1)
				if app.ID != nil {
					fmt.Printf("   🆔 ID   : %s\n", *app.ID)
				}
				if app.Name != nil {
					fmt.Printf("   📛 Name : %s\n", *app.Name)
				}
				if app.Href != nil {
					fmt.Printf("   🔗 Href : %s\n", *app.Href)
				}
			}
		})
	})

	Describe(`CreateApplication - Add an application to an instance`, func() {
		BeforeEach(func() {
			SkipUnlessConfigured(config, "SERVICE_INSTANCE_DEPLOYMENT_GUID")
		})
		It(`CreateApplication(createApplicationOptions *CreateApplicationOptions)`, func() {
			createApplicationOptions := &mqcloudv1.CreateApplicationOptions{
				ServiceInstanceGuid: core.StringPtr(config["SERVICE_INSTANCE_DEPLOYMENT_GUID"]),
				Name:                core.StringPtr("testing123"),
			}

			applicationCreated, response, err := mqcloudService.CreateApplication(createApplicationOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(201))
			Expect(applicationCreated).ToNot(BeNil())

			// Save the ApplicationID to config for use in GetApplication
			config["APPLICATION_ID"] = *applicationCreated.ID

			fmt.Println("✅ Application has been created successfully")

			// Optional pause to allow backend propagation
			time.Sleep(30 * time.Second)
		})
	})

	Describe(`GetApplication - Get an application for an instance`, func() {
		BeforeEach(func() {
			SkipUnlessConfigured(config, "SERVICE_INSTANCE_DEPLOYMENT_GUID", "APPLICATION_ID")
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

	Describe(`SetApplicationName - Update the name for an application`, func() {
		BeforeEach(func() {
			SkipUnlessConfigured(config, "SERVICE_INSTANCE_DEPLOYMENT_GUID", "APPLICATION_ID")
		})
		It(`SetApplicationName(setApplicationNameOptions *SetApplicationNameOptions)`, func() {
			setApplicationNameOptions := &mqcloudv1.SetApplicationNameOptions{
				ServiceInstanceGuid: core.StringPtr(config["SERVICE_INSTANCE_DEPLOYMENT_GUID"]),
				ApplicationID:       core.StringPtr(config["APPLICATION_ID"]),
				Name:                core.StringPtr("test456"),
			}

			applicationDetails, response, err := mqcloudService.SetApplicationName(setApplicationNameOptions)
			Expect(err).To(BeNil())
			Expect(response).ToNot(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(applicationDetails).ToNot(BeNil())

			fmt.Printf("✅ Application shortname has been successfully updated to '%s'\n", *setApplicationNameOptions.Name)
		})
	})

	Describe(`CreateApplicationApikey - Create a new apikey for an application`, func() {
		BeforeEach(func() {
			SkipUnlessConfigured(config, "SERVICE_INSTANCE_DEPLOYMENT_GUID")
		})
		It(`CreateApplicationApikey(createApplicationApikeyOptions *CreateApplicationApikeyOptions)`, func() {
			createApplicationApikeyOptions := &mqcloudv1.CreateApplicationApikeyOptions{
				ServiceInstanceGuid: core.StringPtr(config["SERVICE_INSTANCE_DEPLOYMENT_GUID"]),
				ApplicationID:       core.StringPtr(config["APPLICATION_ID"]),
				Name:                core.StringPtr("testapikey1"),
			}

			applicationApiKeyCreated, response, err := mqcloudService.CreateApplicationApikey(createApplicationApikeyOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(201))
			Expect(applicationApiKeyCreated).ToNot(BeNil())
		})
	})

	Describe(`CreateTrustStorePemCertificate - Upload a trust store certificate`, func() {
		BeforeEach(func() {
			SkipTestIfQmIsNotRunning(config, mqcloudService)
			SkipUnlessConfigured(config, "TRUST_STORE_FILE_PATH")
		})
		It(`CreateTrustStorePemCertificate(createTrustStorePemCertificateOptions *CreateTrustStorePemCertificateOptions)`, func() {
			file, err := os.Open(config["TRUST_STORE_FILE_PATH"])
			if err != nil {
				_, _ = fmt.Fprintf(GinkgoWriter, "Error opening file: %s \n", err.Error())
				return
			}
			defer func() { _ = file.Close() }()
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
			SkipUnlessConfigured(config, "SERVICE_INSTANCE_DEPLOYMENT_GUID", "QUEUE_MANAGER_ID")
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
			SkipUnlessConfigured(config, "SERVICE_INSTANCE_DEPLOYMENT_GUID", "QUEUE_MANAGER_ID", "TRUST_STORE_CERTIFICATE_ID")
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
			SkipUnlessConfigured(config, "SERVICE_INSTANCE_DEPLOYMENT_GUID", "QUEUE_MANAGER_ID", "TRUST_STORE_CERTIFICATE_ID")
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
			SkipTestIfQmIsNotRunning(config, mqcloudService)
			SkipUnlessConfigured(config, "KEY_STORE_FILE_PATH")
		})
		It(`CreateKeyStorePemCertificate(createKeyStorePemCertificateOptions *CreateKeyStorePemCertificateOptions)`, func() {
			file, err := os.Open(config["KEY_STORE_FILE_PATH"])
			if err != nil {
				_, _ = fmt.Fprintf(GinkgoWriter, "Error opening file: %s \n", err.Error())
				return
			}
			defer func() { _ = file.Close() }()
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
			SkipUnlessConfigured(config, "SERVICE_INSTANCE_DEPLOYMENT_GUID", "QUEUE_MANAGER_ID")
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
			SkipUnlessConfigured(config, "SERVICE_INSTANCE_DEPLOYMENT_GUID", "QUEUE_MANAGER_ID", "KEY_STORE_CERTIFICATE_ID")
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
			SkipUnlessConfigured(config, "SERVICE_INSTANCE_DEPLOYMENT_GUID", "QUEUE_MANAGER_ID", "KEY_STORE_CERTIFICATE_ID")
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
			SkipUnlessConfigured(config, "SERVICE_INSTANCE_DEPLOYMENT_GUID", "QUEUE_MANAGER_ID", "KEY_STORE_CERTIFICATE_ID")
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
			SkipUnlessConfigured(config, "SERVICE_INSTANCE_DEPLOYMENT_GUID", "QUEUE_MANAGER_ID", "KEY_STORE_CERTIFICATE_ID")
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
			SkipUnlessConfigured(config, "SERVICE_INSTANCE_CAPACITY_GUID", "TARGET_CRN")
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
			SkipUnlessConfigured(config, "SERVICE_INSTANCE_CAPACITY_GUID")
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
			}
			config["NUMBER_OF_VIRTUAL_GATEWAYS"] = strconv.Itoa(len(allResults))
			_, _ = fmt.Fprintf(GinkgoWriter, "Retrieved a total of %d item(s) with pagination.\n", len(allResults))
		})
		It(`ListVirtualPrivateEndpointGateways(listVirtualPrivateEndpointGatewaysOptions *ListVirtualPrivateEndpointGatewaysOptions) using VirtualPrivateEndpointGatewaysPager`, func() {
			listVirtualPrivateEndpointGatewaysOptions := &mqcloudv1.ListVirtualPrivateEndpointGatewaysOptions{
				ServiceInstanceGuid: core.StringPtr(config["SERVICE_INSTANCE_CAPACITY_GUID"]),
				TrustedProfile:      core.StringPtr(config["TRUSTED_PROFILE"]),
				Limit:               core.Int64Ptr(int64(10)),
			}

			listVirtualPrivateEndpointGatewaysOptions.TrustedProfile = nil
			listVirtualPrivateEndpointGatewaysOptions.Start = nil

			// Test GetAll(), which fetches every page itself.
			pager, err := mqcloudService.NewVirtualPrivateEndpointGatewaysPager(listVirtualPrivateEndpointGatewaysOptions)
			Expect(err).To(BeNil())
			Expect(pager).ToNot(BeNil())

			allItems, err := pager.GetAll()
			Expect(err).To(BeNil())
			//Expect(allItems).ToNot(BeNil())
			totalVirtualGateways, err := strconv.Atoi(config["NUMBER_OF_VIRTUAL_GATEWAYS"])
			if err != nil {
				_, _ = fmt.Fprintf(GinkgoWriter, "Warning failed to convert config[\"NUMBER_OF_VIRTUAL_GATEWAYS\"] (%s) to integer as such default to 0 .\n", config["NUMBER_OF_VIRTUAL_GATEWAYS"])
				totalVirtualGateways = 0
			}
			Expect(len(allItems)).To(Equal(totalVirtualGateways))
		})
	})

	Describe(`GetVirtualPrivateEndpointGateway - Display the information for a specific virtual private endpoint gateway`, func() {
		BeforeEach(func() {
			SkipUnlessConfigured(config, "SERVICE_INSTANCE_CAPACITY_GUID", "VIRTUAL_PRIVATE_ENDPOINT_GATEWAY_GUID")
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

	Describe(`CreateQueueManagerWebConsoleConfiguration - Create a new web console configuration`, func() {
		BeforeEach(func() {
			SkipTestIfQmIsNotRunning(config, mqcloudService)
		})
		It(`CreateQueueManagerWebConsoleConfiguration(createQueueManagerWebConsoleConfigurationOptions *CreateQueueManagerWebConsoleConfigurationOptions)`, func() {
			webConsoleSettingsModel := &mqcloudv1.WebConsoleSettings{
				Certificate: core.StringPtr("provider_managed"),
			}

			createQueueManagerWebConsoleConfigurationOptions := &mqcloudv1.CreateQueueManagerWebConsoleConfigurationOptions{
				ServiceInstanceGuid: core.StringPtr(config["SERVICE_INSTANCE_DEPLOYMENT_GUID"]),
				QueueManagerID:      core.StringPtr(config["QUEUE_MANAGER_ID"]),
				WebConsole:          webConsoleSettingsModel,
			}

			operationStatus, response, err := mqcloudService.CreateQueueManagerWebConsoleConfiguration(createQueueManagerWebConsoleConfigurationOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(202))
			Expect(operationStatus).ToNot(BeNil())

			config["OPERATION_STATUS_ID"] = *operationStatus.ID
			config["QUEUE_MANAGER_WEBCONSOLE_CONFIGURATION_ID"] = *operationStatus.Resource.ID

			err = WaitForOperation(mqcloudService, config["SERVICE_INSTANCE_DEPLOYMENT_GUID"], config["OPERATION_STATUS_ID"])
			Expect(err).To(BeNil())
		})
	})

	Describe(`GetOperation - Get operation status`, func() {
		BeforeEach(func() {
			SkipUnlessConfigured(config, "SERVICE_INSTANCE_DEPLOYMENT_GUID", "OPERATION_STATUS_ID")
		})
		It(`GetOperation(getOperationOptions *GetOperationOptions)`, func() {
			getOperationOptions := &mqcloudv1.GetOperationOptions{
				ServiceInstanceGuid: core.StringPtr(config["SERVICE_INSTANCE_DEPLOYMENT_GUID"]),
				OperationID:         core.StringPtr(config["OPERATION_STATUS_ID"]),
			}

			operationStatus, response, err := mqcloudService.GetOperation(getOperationOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(operationStatus).ToNot(BeNil())
		})
	})

	Describe(`ListQueueManagerWebConsoleConfigurations - Get list of web console configurations`, func() {
		BeforeEach(func() {
			SkipTestIfQmIsNotRunning(config, mqcloudService)
		})
		It(`ListQueueManagerWebConsoleConfigurations(listQueueManagerWebConsoleConfigurationsOptions *ListQueueManagerWebConsoleConfigurationsOptions) with pagination`, func() {
			listQueueManagerWebConsoleConfigurationsOptions := &mqcloudv1.ListQueueManagerWebConsoleConfigurationsOptions{
				ServiceInstanceGuid: core.StringPtr(config["SERVICE_INSTANCE_DEPLOYMENT_GUID"]),
				QueueManagerID:      core.StringPtr(config["QUEUE_MANAGER_ID"]),
			}

			listQueueManagerWebConsoleConfigurationsOptions.Offset = nil
			listQueueManagerWebConsoleConfigurationsOptions.Limit = core.Int64Ptr(1)

			var allResults []mqcloudv1.WebConsoleConfiguration
			for {
				webConsoleConfigurationCollection, response, err := mqcloudService.ListQueueManagerWebConsoleConfigurations(listQueueManagerWebConsoleConfigurationsOptions)
				Expect(err).To(BeNil())
				Expect(response.StatusCode).To(Equal(200))
				Expect(webConsoleConfigurationCollection).ToNot(BeNil())
				allResults = append(allResults, webConsoleConfigurationCollection.WebConsoleConfigurations...)

				nextToken, nextErr := webConsoleConfigurationCollection.GetNextOffset()
				Expect(nextErr).To(BeNil())
				listQueueManagerWebConsoleConfigurationsOptions.Offset = nextToken

				if listQueueManagerWebConsoleConfigurationsOptions.Offset == nil {
					break
				}
			}
			_, _ = fmt.Fprintf(GinkgoWriter, "Retrieved a total of %d item(s) with pagination.\n", len(allResults))
		})
		It(`ListQueueManagerWebConsoleConfigurations(listQueueManagerWebConsoleConfigurationsOptions *ListQueueManagerWebConsoleConfigurationsOptions) using QueueManagerWebConsoleConfigurationsPager`, func() {
			listQueueManagerWebConsoleConfigurationsOptions := &mqcloudv1.ListQueueManagerWebConsoleConfigurationsOptions{
				ServiceInstanceGuid: core.StringPtr(config["SERVICE_INSTANCE_DEPLOYMENT_GUID"]),
				QueueManagerID:      core.StringPtr(config["QUEUE_MANAGER_ID"]),
				Limit:               core.Int64Ptr(int64(10)),
			}

			// Test GetNext().
			pager, err := mqcloudService.NewQueueManagerWebConsoleConfigurationsPager(listQueueManagerWebConsoleConfigurationsOptions)
			Expect(err).To(BeNil())
			Expect(pager).ToNot(BeNil())

			var allResults []mqcloudv1.WebConsoleConfiguration
			for pager.HasNext() {
				nextPage, nextPageErr := pager.GetNext()
				Expect(nextPageErr).To(BeNil())
				Expect(nextPage).ToNot(BeNil())
				allResults = append(allResults, nextPage...)
			}

			// Test GetAll().
			pager, err = mqcloudService.NewQueueManagerWebConsoleConfigurationsPager(listQueueManagerWebConsoleConfigurationsOptions)
			Expect(err).To(BeNil())
			Expect(pager).ToNot(BeNil())

			allItems, err := pager.GetAll()
			Expect(err).To(BeNil())
			Expect(allItems).ToNot(BeNil())

			Expect(len(allItems)).To(Equal(len(allResults)))
			_, _ = fmt.Fprintf(GinkgoWriter, "ListQueueManagerWebConsoleConfigurations() returned a total of %d item(s) using QueueManagerWebConsoleConfigurationsPager.\n", len(allResults))
		})
	})

	Describe(`GetQueueManagerWebConsoleConfiguration - Get details of a web console configuration`, func() {
		BeforeEach(func() {
			SkipTestIfQmIsNotRunning(config, mqcloudService)
			SkipUnlessConfigured(config, "QUEUE_MANAGER_WEBCONSOLE_CONFIGURATION_ID")
		})
		It(`GetQueueManagerWebConsoleConfiguration(getQueueManagerWebConsoleConfigurationOptions *GetQueueManagerWebConsoleConfigurationOptions)`, func() {
			getQueueManagerWebConsoleConfigurationOptions := &mqcloudv1.GetQueueManagerWebConsoleConfigurationOptions{
				ServiceInstanceGuid: core.StringPtr(config["SERVICE_INSTANCE_DEPLOYMENT_GUID"]),
				QueueManagerID:      core.StringPtr(config["QUEUE_MANAGER_ID"]),
				ConfigurationID:     core.StringPtr(config["QUEUE_MANAGER_WEBCONSOLE_CONFIGURATION_ID"]),
			}

			webConsoleConfiguration, response, err := mqcloudService.GetQueueManagerWebConsoleConfiguration(getQueueManagerWebConsoleConfigurationOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(webConsoleConfiguration).ToNot(BeNil())
		})
	})

	Describe(`ReplaceQueueManagerWebConsoleConfiguration - Replace a web console configuration`, func() {
		BeforeEach(func() {
			SkipTestIfQmIsNotRunning(config, mqcloudService)
			SkipUnlessConfigured(config, "QUEUE_MANAGER_WEBCONSOLE_CONFIGURATION_ID")
		})
		It(`ReplaceQueueManagerWebConsoleConfiguration(replaceQueueManagerWebConsoleConfigurationOptions *ReplaceQueueManagerWebConsoleConfigurationOptions)`, func() {
			webConsoleSettingsModel := &mqcloudv1.WebConsoleSettings{
				Certificate: core.StringPtr("provider_managed"),
			}

			replaceQueueManagerWebConsoleConfigurationOptions := &mqcloudv1.ReplaceQueueManagerWebConsoleConfigurationOptions{
				ServiceInstanceGuid: core.StringPtr(config["SERVICE_INSTANCE_DEPLOYMENT_GUID"]),
				QueueManagerID:      core.StringPtr(config["QUEUE_MANAGER_ID"]),
				ConfigurationID:     core.StringPtr(config["QUEUE_MANAGER_WEBCONSOLE_CONFIGURATION_ID"]),
				WebConsole:          webConsoleSettingsModel,
			}

			operationStatus, response, err := mqcloudService.ReplaceQueueManagerWebConsoleConfiguration(replaceQueueManagerWebConsoleConfigurationOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(202))
			Expect(operationStatus).ToNot(BeNil())

			err = WaitForOperation(mqcloudService, config["SERVICE_INSTANCE_DEPLOYMENT_GUID"], *operationStatus.ID)
			Expect(err).To(BeNil())
		})
	})

	Describe(`CreateQueueManagerConfiguration - Create a new queue manager configuration`, func() {
		BeforeEach(func() {
			SkipTestIfQmIsNotRunning(config, mqcloudService)
		})
		It(`CreateQueueManagerConfiguration(createQueueManagerConfigurationOptions *CreateQueueManagerConfigurationOptions)`, func() {
			queueManagerSettingsModel := &mqcloudv1.QueueManagerSettings{
				Certificate: core.StringPtr("provider_managed"),
			}

			createQueueManagerConfigurationOptions := &mqcloudv1.CreateQueueManagerConfigurationOptions{
				ServiceInstanceGuid: core.StringPtr(config["SERVICE_INSTANCE_DEPLOYMENT_GUID"]),
				QueueManagerID:      core.StringPtr(config["QUEUE_MANAGER_ID"]),
				QueueManager:        queueManagerSettingsModel,
			}

			operationStatus, response, err := mqcloudService.CreateQueueManagerConfiguration(createQueueManagerConfigurationOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(202))
			Expect(operationStatus).ToNot(BeNil())

			config["OPERATION_STATUS_ID"] = *operationStatus.ID
			config["QUEUE_MANAGER_CONFIGURATION_ID"] = *operationStatus.Resource.ID

			err = WaitForOperation(mqcloudService, config["SERVICE_INSTANCE_DEPLOYMENT_GUID"], config["OPERATION_STATUS_ID"])
			Expect(err).To(BeNil())
		})
	})

	Describe(`ListQueueManagerConfigurations - Get list of queue manager configurations`, func() {
		BeforeEach(func() {
			SkipTestIfQmIsNotRunning(config, mqcloudService)
		})
		It(`ListQueueManagerConfigurations(listQueueManagerConfigurationsOptions *ListQueueManagerConfigurationsOptions) with pagination`, func() {
			listQueueManagerConfigurationsOptions := &mqcloudv1.ListQueueManagerConfigurationsOptions{
				ServiceInstanceGuid: core.StringPtr(config["SERVICE_INSTANCE_DEPLOYMENT_GUID"]),
				QueueManagerID:      core.StringPtr(config["QUEUE_MANAGER_ID"]),
			}

			listQueueManagerConfigurationsOptions.Offset = nil
			listQueueManagerConfigurationsOptions.Limit = core.Int64Ptr(1)

			var allResults []mqcloudv1.QueueManagerConfiguration
			for {
				queueManagerConfigurationCollection, response, err := mqcloudService.ListQueueManagerConfigurations(listQueueManagerConfigurationsOptions)
				Expect(err).To(BeNil())
				Expect(response.StatusCode).To(Equal(200))
				Expect(queueManagerConfigurationCollection).ToNot(BeNil())
				allResults = append(allResults, queueManagerConfigurationCollection.QueueManagerConfigurations...)

				nextToken, nextErr := queueManagerConfigurationCollection.GetNextOffset()
				Expect(nextErr).To(BeNil())
				listQueueManagerConfigurationsOptions.Offset = nextToken

				if listQueueManagerConfigurationsOptions.Offset == nil {
					break
				}
			}
			_, _ = fmt.Fprintf(GinkgoWriter, "Retrieved a total of %d item(s) with pagination.\n", len(allResults))
		})
		It(`ListQueueManagerConfigurations(listQueueManagerConfigurationsOptions *ListQueueManagerConfigurationsOptions) using QueueManagerConfigurationsPager`, func() {
			listQueueManagerConfigurationsOptions := &mqcloudv1.ListQueueManagerConfigurationsOptions{
				ServiceInstanceGuid: core.StringPtr(config["SERVICE_INSTANCE_DEPLOYMENT_GUID"]),
				QueueManagerID:      core.StringPtr(config["QUEUE_MANAGER_ID"]),
				Limit:               core.Int64Ptr(int64(10)),
			}

			// Test GetNext().
			pager, err := mqcloudService.NewQueueManagerConfigurationsPager(listQueueManagerConfigurationsOptions)
			Expect(err).To(BeNil())
			Expect(pager).ToNot(BeNil())

			var allResults []mqcloudv1.QueueManagerConfiguration
			for pager.HasNext() {
				nextPage, nextPageErr := pager.GetNext()
				Expect(nextPageErr).To(BeNil())
				Expect(nextPage).ToNot(BeNil())
				allResults = append(allResults, nextPage...)
			}

			// Test GetAll().
			pager, err = mqcloudService.NewQueueManagerConfigurationsPager(listQueueManagerConfigurationsOptions)
			Expect(err).To(BeNil())
			Expect(pager).ToNot(BeNil())

			allItems, err := pager.GetAll()
			Expect(err).To(BeNil())
			Expect(allItems).ToNot(BeNil())

			Expect(len(allItems)).To(Equal(len(allResults)))
			_, _ = fmt.Fprintf(GinkgoWriter, "ListQueueManagerConfigurations() returned a total of %d item(s) using QueueManagerConfigurationsPager.\n", len(allResults))
		})
	})

	Describe(`GetQueueManagerConfiguration - Get details of a queue manager configuration`, func() {
		BeforeEach(func() {
			SkipTestIfQmIsNotRunning(config, mqcloudService)
			SkipUnlessConfigured(config, "QUEUE_MANAGER_CONFIGURATION_ID")
		})
		It(`GetQueueManagerConfiguration(getQueueManagerConfigurationOptions *GetQueueManagerConfigurationOptions)`, func() {
			getQueueManagerConfigurationOptions := &mqcloudv1.GetQueueManagerConfigurationOptions{
				ServiceInstanceGuid: core.StringPtr(config["SERVICE_INSTANCE_DEPLOYMENT_GUID"]),
				QueueManagerID:      core.StringPtr(config["QUEUE_MANAGER_ID"]),
				ConfigurationID:     core.StringPtr(config["QUEUE_MANAGER_CONFIGURATION_ID"]),
			}

			queueManagerConfiguration, response, err := mqcloudService.GetQueueManagerConfiguration(getQueueManagerConfigurationOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(queueManagerConfiguration).ToNot(BeNil())
		})
	})

	Describe(`ReplaceQueueManagerConfiguration - Replace a queue manager configuration`, func() {
		BeforeEach(func() {
			SkipTestIfQmIsNotRunning(config, mqcloudService)
			SkipUnlessConfigured(config, "QUEUE_MANAGER_CONFIGURATION_ID")
		})
		It(`ReplaceQueueManagerConfiguration(replaceQueueManagerConfigurationOptions *ReplaceQueueManagerConfigurationOptions)`, func() {
			queueManagerSettingsModel := &mqcloudv1.QueueManagerSettings{
				Certificate: core.StringPtr("provider_managed"),
			}

			replaceQueueManagerConfigurationOptions := &mqcloudv1.ReplaceQueueManagerConfigurationOptions{
				ServiceInstanceGuid: core.StringPtr(config["SERVICE_INSTANCE_DEPLOYMENT_GUID"]),
				QueueManagerID:      core.StringPtr(config["QUEUE_MANAGER_ID"]),
				ConfigurationID:     core.StringPtr(config["QUEUE_MANAGER_CONFIGURATION_ID"]),
				QueueManager:        queueManagerSettingsModel,
			}

			operationStatus, response, err := mqcloudService.ReplaceQueueManagerConfiguration(replaceQueueManagerConfigurationOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(202))
			Expect(operationStatus).ToNot(BeNil())

			err = WaitForOperation(mqcloudService, config["SERVICE_INSTANCE_DEPLOYMENT_GUID"], *operationStatus.ID)
			Expect(err).To(BeNil())
		})
	})

	Describe(`DeleteQueueManagerWebConsoleConfiguration - Delete a web console configuration`, func() {
		BeforeEach(func() {
			SkipTestIfQmIsNotRunning(config, mqcloudService)
			SkipUnlessConfigured(config, "QUEUE_MANAGER_WEBCONSOLE_CONFIGURATION_ID")
		})
		It(`DeleteQueueManagerWebConsoleConfiguration(deleteQueueManagerWebConsoleConfigurationOptions *DeleteQueueManagerWebConsoleConfigurationOptions)`, func() {
			deleteQueueManagerWebConsoleConfigurationOptions := &mqcloudv1.DeleteQueueManagerWebConsoleConfigurationOptions{
				ServiceInstanceGuid: core.StringPtr(config["SERVICE_INSTANCE_DEPLOYMENT_GUID"]),
				QueueManagerID:      core.StringPtr(config["QUEUE_MANAGER_ID"]),
				ConfigurationID:     core.StringPtr(config["QUEUE_MANAGER_WEBCONSOLE_CONFIGURATION_ID"]),
			}

			operationStatus, response, err := mqcloudService.DeleteQueueManagerWebConsoleConfiguration(deleteQueueManagerWebConsoleConfigurationOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(202))
			Expect(operationStatus).ToNot(BeNil())

			err = WaitForOperation(mqcloudService, config["SERVICE_INSTANCE_DEPLOYMENT_GUID"], *operationStatus.ID)
			Expect(err).To(BeNil())
		})
	})

	Describe(`DeleteQueueManagerConfiguration - Delete a queue manager configuration`, func() {
		BeforeEach(func() {
			SkipTestIfQmIsNotRunning(config, mqcloudService)
			SkipUnlessConfigured(config, "QUEUE_MANAGER_CONFIGURATION_ID")
		})
		It(`DeleteQueueManagerConfiguration(deleteQueueManagerConfigurationOptions *DeleteQueueManagerConfigurationOptions)`, func() {
			deleteQueueManagerConfigurationOptions := &mqcloudv1.DeleteQueueManagerConfigurationOptions{
				ServiceInstanceGuid: core.StringPtr(config["SERVICE_INSTANCE_DEPLOYMENT_GUID"]),
				QueueManagerID:      core.StringPtr(config["QUEUE_MANAGER_ID"]),
				ConfigurationID:     core.StringPtr(config["QUEUE_MANAGER_CONFIGURATION_ID"]),
			}

			operationStatus, response, err := mqcloudService.DeleteQueueManagerConfiguration(deleteQueueManagerConfigurationOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(202))
			Expect(operationStatus).ToNot(BeNil())

			err = WaitForOperation(mqcloudService, config["SERVICE_INSTANCE_DEPLOYMENT_GUID"], *operationStatus.ID)
			Expect(err).To(BeNil())
		})
	})

	Describe(`DeleteKeyStoreCertificate - Delete a queue manager's key store certificate`, func() {
		BeforeEach(func() {
			SkipUnlessConfigured(config, "SERVICE_INSTANCE_DEPLOYMENT_GUID", "QUEUE_MANAGER_ID", "KEY_STORE_CERTIFICATE_ID")
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

	Describe(`DeleteTrustStoreCertificate - Delete a trust store certificate`, func() {
		BeforeEach(func() {
			SkipUnlessConfigured(config, "SERVICE_INSTANCE_DEPLOYMENT_GUID", "QUEUE_MANAGER_ID", "TRUST_STORE_CERTIFICATE_ID")
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

	Describe(`DeleteVirtualPrivateEndpointGateway - Delete a specific virtual private endpoint gateway`, func() {
		BeforeEach(func() {
			SkipUnlessConfigured(config, "SERVICE_INSTANCE_CAPACITY_GUID", "VIRTUAL_PRIVATE_ENDPOINT_GATEWAY_GUID")
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

	Describe(`DeleteUser - Delete a user for an instance`, func() {
		BeforeEach(func() {
			SkipUnlessConfigured(config, "SERVICE_INSTANCE_DEPLOYMENT_GUID", "USER_ID")
		})
		It(`DeleteUser(deleteUserOptions *DeleteUserOptions)`, func() {
			deleteUserOptions := &mqcloudv1.DeleteUserOptions{
				ServiceInstanceGuid: core.StringPtr(config["SERVICE_INSTANCE_DEPLOYMENT_GUID"]),
				UserID:              core.StringPtr(config["USER_ID"]),
			}

			response, err := mqcloudService.DeleteUser(deleteUserOptions)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(204))
			fmt.Println("✅ User has been deleted successfully")
		})
	})

	Describe(`DeleteApplication - Delete an application from an instance`, func() {
		BeforeEach(func() {
			SkipUnlessConfigured(config, "SERVICE_INSTANCE_DEPLOYMENT_GUID", "APPLICATION_ID")
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

	Describe(`DeleteQueueManager - Delete a queue manager`, func() {
		BeforeEach(func() {
			SkipUnlessConfigured(config, "SERVICE_INSTANCE_DEPLOYMENT_GUID", "QUEUE_MANAGER_ID")
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
})
