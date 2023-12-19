//go:build examples
// +build examples

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
	"encoding/json"
	"fmt"
	//	"io"
	"os"

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
		serviceURL     string

		queue_manager_id          string
		// user_id                   string
		// application_id            string
		truststore_certificate_id string
		// keystore_certificate_id   string
		serviceinstance_guid      string
		// keystore_filepath         string
		truststore_filepath       string
		queue_manager_id_ref      *string
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

			serviceURL = os.Getenv("IBMCLOUD_MQCLOUD_CONFIG_ENDPOINT")
			if serviceURL == "" {
				Skip("Unable to load service URL configuration property, skipping tests")
			}
			serviceinstance_guid = os.Getenv("SERVICEINSTANCE_GUID")
			// keystore_filepath = os.Getenv("KEYSTORE_FILE_PATH")
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
			var err error

			// begin-common
			apikey := os.Getenv("IC_API_KEY")
			authenticator := &core.IamAuthenticator{
				ApiKey: apikey,
				URL:    os.Getenv("IBMCLOUD_IAM_API_ENDPOINT") + "/identity/token",
			}
			acceptLanguage := "en-US"
			if err != nil {
				panic(err)
			}
			mqcloudServiceOptions := &mqcloudv1.MqcloudV1Options{
				URL:            os.Getenv("IBMCLOUD_MQCLOUD_CONFIG_ENDPOINT"),
				Authenticator:  authenticator,
				ServiceName:    mqcloudv1.DefaultServiceName,
				AcceptLanguage: core.StringPtr(acceptLanguage),
			}

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
		/*	It(`GetUsageDetails request example`, func() {
				fmt.Println("\nGetUsageDetails() result:")
				// begin-get_usage_details

				getUsageDetailsOptions := mqcloudService.NewGetUsageDetailsOptions(
					serviceinstance_guid,
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
					serviceinstance_guid,
				)

				configurationOptions, response, err := mqcloudService.GetOptions(getOptionsOptions)
				if err != nil {
					panic(err)
				}
				b, _ := json.MarshalIndent(configurationOptions, "", "  ")
				fmt.Println(string(b))

				// end-get_options

				Expect(err).To(BeNil())
				Expect(response.StatusCode).To(Equal(200))
				Expect(configurationOptions).ToNot(BeNil())
			})*/
		It(`CreateQueueManager request example`, func() {
			fmt.Println("\nCreateQueueManager() result:")
			// begin-create_queue_manager

			createQueueManagerOptions := mqcloudService.NewCreateQueueManagerOptions(
				serviceinstance_guid,
				"testqm"+RandString(6),
				"ibmcloud_eu_de",
				"lite",
			)

			queueManagerTaskStatus, response, err := mqcloudService.CreateQueueManager(createQueueManagerOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(queueManagerTaskStatus, "", "  ")
			fmt.Println(string(b))

			// end-create_queue_manager

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(202))
			Expect(queueManagerTaskStatus).ToNot(BeNil())
			queue_manager_id = *queueManagerTaskStatus.QueueManagerID
			queue_manager_id_ref = queueManagerTaskStatus.QueueManagerID
		})
		// It(`ListQueueManagers request example`, func() {
		// 	fmt.Println("\nListQueueManagers() result:")
		// 	// begin-list_queue_managers
		// 	listQueueManagersOptions := &mqcloudv1.ListQueueManagersOptions{
		// 		ServiceInstanceGuid: core.StringPtr(serviceinstance_guid),
		// 		Limit:               core.Int64Ptr(int64(10)),
		// 	}

		// 	pager, err := mqcloudService.NewQueueManagersPager(listQueueManagersOptions)
		// 	if err != nil {
		// 		panic(err)
		// 	}

		// 	var allResults []mqcloudv1.QueueManagerDetails
		// 	for pager.HasNext() {
		// 		nextPage, err := pager.GetNext()
		// 		if err != nil {
		// 			panic(err)
		// 		}
		// 		allResults = append(allResults, nextPage...)
		// 	}
		// 	b, _ := json.MarshalIndent(allResults, "", "  ")
		// 	fmt.Println(string(b))
		// 	// end-list_queue_managers
		// })
		It(`GetQueueManager request example`, func() {
			fmt.Println("\nGetQueueManager() result:")
			// begin-get_queue_manager

			getQueueManagerOptions := mqcloudService.NewGetQueueManagerOptions(
				serviceinstance_guid,
				queue_manager_id,
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
		/*	It(`SetQueueManagerVersion request example`, func() {
				fmt.Println("\nSetQueueManagerVersion() result:")
				// begin-set_queue_manager_version
				error := WaitForQmStatusUpdate(queue_manager_id_ref, mqcloudService, serviceinstance_guid)
			 if error != nil {
				 fmt.Fprintf(GinkgoWriter, "WaitForQmStatusUpdate failed: %s \n", error)
				 return
			 }
			 fmt.Fprintf(GinkgoWriter,
				 "--------- Queue Manager is now in the running state ---------",
			 )

				setQueueManagerVersionOptions := mqcloudService.NewSetQueueManagerVersionOptions(
					serviceinstance_guid,
					queue_manager_id,
					"9.3.4_3",
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
			})*/
		// It(`GetQueueManagerAvailableUpgradeVersions request example`, func() {
		// 	fmt.Println("\nGetQueueManagerAvailableUpgradeVersions() result:")
		// 	// begin-get_queue_manager_available_upgrade_versions
		// 	error := WaitForQmStatusUpdate(queue_manager_id_ref, mqcloudService, serviceinstance_guid)
		// 	if error != nil {
		// 		fmt.Fprintf(GinkgoWriter, "WaitForQmStatusUpdate failed: %s \n", error)
		// 		return
		// 	}
		// 	fmt.Fprintf(GinkgoWriter,
		// 		"--------- Queue Manager is now in the running state ---------",
		// 	)

		// 	getQueueManagerAvailableUpgradeVersionsOptions := mqcloudService.NewGetQueueManagerAvailableUpgradeVersionsOptions(
		// 		serviceinstance_guid,
		// 		queue_manager_id,
		// 	)

		// 	queueManagerVersionUpgrades, response, err := mqcloudService.GetQueueManagerAvailableUpgradeVersions(getQueueManagerAvailableUpgradeVersionsOptions)
		// 	if err != nil {
		// 		panic(err)
		// 	}
		// 	b, _ := json.MarshalIndent(queueManagerVersionUpgrades, "", "  ")
		// 	fmt.Println(string(b))

		// 	// end-get_queue_manager_available_upgrade_versions

		// 	Expect(err).To(BeNil())
		// 	Expect(response.StatusCode).To(Equal(200))
		// 	Expect(queueManagerVersionUpgrades).ToNot(BeNil())
		// })
		// It(`GetQueueManagerConnectionInfo request example`, func() {
		// 	fmt.Println("\nGetQueueManagerConnectionInfo() result:")
		// 	// begin-get_queue_manager_connection_info

		// 	getQueueManagerConnectionInfoOptions := mqcloudService.NewGetQueueManagerConnectionInfoOptions(
		// 		serviceinstance_guid,
		// 		queue_manager_id,
		// 	)

		// 	connectionInfo, response, err := mqcloudService.GetQueueManagerConnectionInfo(getQueueManagerConnectionInfoOptions)
		// 	if err != nil {
		// 		panic(err)
		// 	}
		// 	b, _ := json.MarshalIndent(connectionInfo, "", "  ")
		// 	fmt.Println(string(b))

		// 	// end-get_queue_manager_connection_info

		// 	Expect(err).To(BeNil())
		// 	Expect(response.StatusCode).To(Equal(200))
		// 	Expect(connectionInfo).ToNot(BeNil())
		// })
		// It(`GetQueueManagerStatus request example`, func() {
		// 	fmt.Println("\nGetQueueManagerStatus() result:")
		// 	// begin-get_queue_manager_status

		// 	getQueueManagerStatusOptions := mqcloudService.NewGetQueueManagerStatusOptions(
		// 		serviceinstance_guid,
		// 		queue_manager_id,
		// 	)

		// 	queueManagerStatus, response, err := mqcloudService.GetQueueManagerStatus(getQueueManagerStatusOptions)
		// 	if err != nil {
		// 		panic(err)
		// 	}
		// 	b, _ := json.MarshalIndent(queueManagerStatus, "", "  ")
		// 	fmt.Println(string(b))

		// 	// end-get_queue_manager_status

		// 	Expect(err).To(BeNil())
		// 	Expect(response.StatusCode).To(Equal(200))
		// 	Expect(queueManagerStatus).ToNot(BeNil())
		// })
		// It(`ListUsers request example`, func() {
		// 	fmt.Println("\nListUsers() result:")
		// 	// begin-list_users
		// 	listUsersOptions := &mqcloudv1.ListUsersOptions{
		// 		ServiceInstanceGuid: core.StringPtr(serviceinstance_guid),
		// 		Limit:               core.Int64Ptr(int64(10)),
		// 	}

		// 	pager, err := mqcloudService.NewUsersPager(listUsersOptions)
		// 	if err != nil {
		// 		panic(err)
		// 	}

		// 	var allResults []mqcloudv1.UserDetails
		// 	for pager.HasNext() {
		// 		nextPage, err := pager.GetNext()
		// 		if err != nil {
		// 			panic(err)
		// 		}
		// 		allResults = append(allResults, nextPage...)
		// 	}
		// 	b, _ := json.MarshalIndent(allResults, "", "  ")
		// 	fmt.Println(string(b))
		// 	// end-list_users
		// })
		// It(`CreateUser request example`, func() {
		// 	fmt.Println("\nCreateUser() result:")
		// 	// begin-create_user

		// 	createUserOptions := mqcloudService.NewCreateUserOptions(
		// 		serviceinstance_guid,
		// 		"user"+RandString(6)+"@ibm.com",
		// 		"user"+RandString(6),
		// 	)

		// 	userDetails, response, err := mqcloudService.CreateUser(createUserOptions)
		// 	if err != nil {
		// 		panic(err)
		// 	}
		// 	b, _ := json.MarshalIndent(userDetails, "", "  ")
		// 	fmt.Println(string(b))

		// 	// end-create_user

		// 	Expect(err).To(BeNil())
		// 	Expect(response.StatusCode).To(Equal(201))
		// 	Expect(userDetails).ToNot(BeNil())
		// 	user_id = *userDetails.ID
		// })
		// It(`GetUser request example`, func() {
		// 	fmt.Println("\nGetUser() result:")
		// 	// begin-get_user

		// 	getUserOptions := mqcloudService.NewGetUserOptions(
		// 		serviceinstance_guid,
		// 		user_id,
		// 	)

		// 	userDetails, response, err := mqcloudService.GetUser(getUserOptions)
		// 	if err != nil {
		// 		panic(err)
		// 	}
		// 	b, _ := json.MarshalIndent(userDetails, "", "  ")
		// 	fmt.Println(string(b))

		// 	// end-get_user

		// 	Expect(err).To(BeNil())
		// 	Expect(response.StatusCode).To(Equal(200))
		// 	Expect(userDetails).ToNot(BeNil())
		// })
		// It(`ListApplications request example`, func() {
		// 	fmt.Println("\nListApplications() result:")
		// 	// begin-list_applications
		// 	listApplicationsOptions := &mqcloudv1.ListApplicationsOptions{
		// 		ServiceInstanceGuid: core.StringPtr(serviceinstance_guid),
		// 		Limit:               core.Int64Ptr(int64(10)),
		// 	}

		// 	pager, err := mqcloudService.NewApplicationsPager(listApplicationsOptions)
		// 	if err != nil {
		// 		panic(err)
		// 	}

		// 	var allResults []mqcloudv1.ApplicationDetails
		// 	for pager.HasNext() {
		// 		nextPage, err := pager.GetNext()
		// 		if err != nil {
		// 			panic(err)
		// 		}
		// 		allResults = append(allResults, nextPage...)
		// 	}
		// 	b, _ := json.MarshalIndent(allResults, "", "  ")
		// 	fmt.Println(string(b))
		// 	// end-list_applications
		// })
		// It(`CreateApplication request example`, func() {
		// 	fmt.Println("\nCreateApplication() result:")
		// 	// begin-create_application

		// 	createApplicationOptions := mqcloudService.NewCreateApplicationOptions(
		// 		serviceinstance_guid,
		// 		"app"+RandString(1),
		// 	)

		// 	applicationCreated, response, err := mqcloudService.CreateApplication(createApplicationOptions)
		// 	if err != nil {
		// 		panic(err)
		// 	}
		// 	b, _ := json.MarshalIndent(applicationCreated, "", "  ")
		// 	fmt.Println(string(b))

		// 	// end-create_application

		// 	Expect(err).To(BeNil())
		// 	Expect(response.StatusCode).To(Equal(201))
		// 	Expect(applicationCreated).ToNot(BeNil())
		// 	application_id = *applicationCreated.ID
		// })
		// It(`GetApplication request example`, func() {
		// 	fmt.Println("\nGetApplication() result:")
		// 	// begin-get_application

		// 	getApplicationOptions := mqcloudService.NewGetApplicationOptions(
		// 		serviceinstance_guid,
		// 		application_id,
		// 	)

		// 	applicationDetails, response, err := mqcloudService.GetApplication(getApplicationOptions)
		// 	if err != nil {
		// 		panic(err)
		// 	}
		// 	b, _ := json.MarshalIndent(applicationDetails, "", "  ")
		// 	fmt.Println(string(b))

		// 	// end-get_application

		// 	Expect(err).To(BeNil())
		// 	Expect(response.StatusCode).To(Equal(200))
		// 	Expect(applicationDetails).ToNot(BeNil())
		// })
		// It(`CreateApplicationApikey request example`, func() {
		// 	fmt.Println("\nCreateApplicationApikey() result:")
		// 	// begin-create_application_apikey

		// 	createApplicationApikeyOptions := mqcloudService.NewCreateApplicationApikeyOptions(
		// 		serviceinstance_guid,
		// 		application_id,
		// 		"test-api-key",
		// 	)

		// 	applicationApiKeyCreated, response, err := mqcloudService.CreateApplicationApikey(createApplicationApikeyOptions)
		// 	if err != nil {
		// 		panic(err)
		// 	}
		// 	b, _ := json.MarshalIndent(applicationApiKeyCreated, "", "  ")
		// 	fmt.Println(string(b))

		// 	// end-create_application_apikey

		// 	Expect(err).To(BeNil())
		// 	Expect(response.StatusCode).To(Equal(201))
		// 	Expect(applicationApiKeyCreated).ToNot(BeNil())
		// })
		It(`CreateTrustStorePemCertificate request example`, func() {
			error := WaitForQmStatusUpdate(queue_manager_id_ref, mqcloudService, serviceinstance_guid)
			if error != nil {
				fmt.Fprintf(GinkgoWriter, "WaitForQmStatusUpdate failed: %s \n", error)
				return
			}
			fmt.Fprintf(GinkgoWriter,
				"--------- Queue Manager is now in the running state ---------",
			)
			fmt.Println("\nCreateTrustStorePemCertificate() result:")
			// begin-create_trust_store_pem_certificate
			// certBytes, err := os.ReadFile(truststore_filepath) // just pass the file name
			// if err != nil {
			// 	fmt.Print(err)
			// }
			// certString := string(certBytes) // convert content to a 'string'
			// rc := io.NopCloser(strings.NewReader(certString))

			filePath := truststore_filepath // Replace with your file path
			file, err := os.Open(filePath)
			if err != nil {
				fmt.Fprintf(GinkgoWriter, "Error opening file: %s \n", err.Error())
				return
			}
			defer file.Close()

			createTrustStorePemCertificateOptions := mqcloudService.NewCreateTrustStorePemCertificateOptions(
				serviceinstance_guid,
				queue_manager_id,
				"testts",
				file,
			)

			trustStoreCertificateDetails, response, err := mqcloudService.CreateTrustStorePemCertificate(createTrustStorePemCertificateOptions)
			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(trustStoreCertificateDetails, "", "  ")
			fmt.Println(string(b))

			// end-create_trust_store_pem_certificate

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(201))
			Expect(trustStoreCertificateDetails).ToNot(BeNil())
			truststore_certificate_id = *trustStoreCertificateDetails.ID
		})
		It(`ListTrustStoreCertificates request example`, func() {
			fmt.Println("\nListTrustStoreCertificates() result:")
			// begin-list_trust_store_certificates

			listTrustStoreCertificatesOptions := mqcloudService.NewListTrustStoreCertificatesOptions(
				serviceinstance_guid,
				queue_manager_id,
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
				serviceinstance_guid,
				queue_manager_id,
				truststore_certificate_id,
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
		/*	It(`DownloadTrustStoreCertificate request example`, func() {
			fmt.Println("\nDownloadTrustStoreCertificate() result:")
			// begin-download_trust_store_certificate
			downloadTrustStoreCertificateOptions := mqcloudService.NewDownloadTrustStoreCertificateOptions(
				serviceinstance_guid,
				queue_manager_id,
				"9b7d1e723af8233",
			)

			result, response, err := mqcloudService.DownloadTrustStoreCertificate(downloadTrustStoreCertificateOptions)
			if err != nil {
				panic(err)
			}
			if result != nil {
				defer result.Close()
				outFile, err := os.Create("result.out")
				if err != nil { panic(err) }
				defer outFile.Close()
				_, err = io.Copy(outFile, result)
				if err != nil { panic(err) }
			}

			// end-download_trust_store_certificate

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(result).ToNot(BeNil())
		})*/
		// It(`CreateKeyStorePemCertificate request example`, func() {
		// 	fmt.Println("\nCreateKeyStorePemCertificate() result:")
		// 	// begin-create_key_store_pem_certificate
		// 	filePath := keystore_filepath // Replace with your file path
		// 	file, err := os.Open(filePath)
		// 	if err != nil {
		// 		fmt.Fprintf(GinkgoWriter, "Error opening file: %s \n", err.Error())
		// 		return
		// 	}
		// 	defer file.Close()
		// 	createKeyStorePemCertificateOptions := mqcloudService.NewCreateKeyStorePemCertificateOptions(
		// 		serviceinstance_guid,
		// 		queue_manager_id,
		// 		"testks",
		// 		file,
		// 	)

		// 	keyStoreCertificateDetails, response, err := mqcloudService.CreateKeyStorePemCertificate(createKeyStorePemCertificateOptions)
		// 	if err != nil {
		// 		panic(err)
		// 	}
		// 	b, _ := json.MarshalIndent(keyStoreCertificateDetails, "", "  ")
		// 	fmt.Println(string(b))

		// 	// end-create_key_store_pem_certificate

		// 	Expect(err).To(BeNil())
		// 	Expect(response.StatusCode).To(Equal(201))
		// 	Expect(keyStoreCertificateDetails).ToNot(BeNil())
		// 	keystore_certificate_id = *keyStoreCertificateDetails.ID
		// })
		// It(`ListKeyStoreCertificates request example`, func() {
		// 	fmt.Println("\nListKeyStoreCertificates() result:")
		// 	// begin-list_key_store_certificates

		// 	listKeyStoreCertificatesOptions := mqcloudService.NewListKeyStoreCertificatesOptions(
		// 		serviceinstance_guid,
		// 		queue_manager_id,
		// 	)

		// 	keyStoreCertificateDetailsCollection, response, err := mqcloudService.ListKeyStoreCertificates(listKeyStoreCertificatesOptions)
		// 	if err != nil {
		// 		panic(err)
		// 	}
		// 	b, _ := json.MarshalIndent(keyStoreCertificateDetailsCollection, "", "  ")
		// 	fmt.Println(string(b))

		// 	// end-list_key_store_certificates

		// 	Expect(err).To(BeNil())
		// 	Expect(response.StatusCode).To(Equal(200))
		// 	Expect(keyStoreCertificateDetailsCollection).ToNot(BeNil())
		// })
		// It(`GetKeyStoreCertificate request example`, func() {
		// 	fmt.Println("\nGetKeyStoreCertificate() result:")
		// 	// begin-get_key_store_certificate

		// 	getKeyStoreCertificateOptions := mqcloudService.NewGetKeyStoreCertificateOptions(
		// 		serviceinstance_guid,
		// 		queue_manager_id,
		// 		keystore_certificate_id,
		// 	)

		// 	keyStoreCertificateDetails, response, err := mqcloudService.GetKeyStoreCertificate(getKeyStoreCertificateOptions)
		// 	if err != nil {
		// 		panic(err)
		// 	}
		// 	b, _ := json.MarshalIndent(keyStoreCertificateDetails, "", "  ")
		// 	fmt.Println(string(b))

		// 	// end-get_key_store_certificate

		// 	Expect(err).To(BeNil())
		// 	Expect(response.StatusCode).To(Equal(200))
		// 	Expect(keyStoreCertificateDetails).ToNot(BeNil())
		// })
		/*	It(`DownloadKeyStoreCertificate request example`, func() {
			fmt.Println("\nDownloadKeyStoreCertificate() result:")
			// begin-download_key_store_certificate

			downloadKeyStoreCertificateOptions := mqcloudService.NewDownloadKeyStoreCertificateOptions(
				serviceinstance_guid,
				queue_manager_id,
				"9b7d1e723af8233",
			)

			result, response, err := mqcloudService.DownloadKeyStoreCertificate(downloadKeyStoreCertificateOptions)
			if err != nil {
				panic(err)
			}
			if result != nil {
				defer result.Close()
				outFile, err := os.Create("result.out")
				if err != nil { panic(err) }
				defer outFile.Close()
				_, err = io.Copy(outFile, result)
				if err != nil { panic(err) }
			}

			// end-download_key_store_certificate

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(200))
			Expect(result).ToNot(BeNil())
		})*/
		/*	It(`GetCertificateAmsChannels request example`, func() {
				fmt.Println("\nGetCertificateAmsChannels() result:")
				// begin-get_certificate_ams_channels

				getCertificateAmsChannelsOptions := mqcloudService.NewGetCertificateAmsChannelsOptions(
					queue_manager_id,
					keystore_certificate_id,
					serviceinstance_guid,
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

				channelDetailsModel := &mqcloudv1.ChannelDetails{
				}

				setCertificateAmsChannelsOptions := mqcloudService.NewSetCertificateAmsChannelsOptions(
					queue_manager_id,
					keystore_certificate_id,
					serviceinstance_guid,
					[]mqcloudv1.ChannelDetails{*channelDetailsModel},
				)

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
			})*/
		It(`DeleteTrustStoreCertificate request example`, func() {
			// begin-delete_trust_store_certificate
			fmt.Println("\nDeleteTrustStoreCertificate() result:", truststore_certificate_id)
			deleteTrustStoreCertificateOptions := mqcloudService.NewDeleteTrustStoreCertificateOptions(
				serviceinstance_guid,
				queue_manager_id,
				truststore_certificate_id,
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
		// It(`DeleteUser request example`, func() {
		// 	// begin-delete_user
		// 	fmt.Println("\nDeleteUser() result:")
		// 	deleteUserOptions := mqcloudService.NewDeleteUserOptions(
		// 		serviceinstance_guid,
		// 		user_id,
		// 	)

		// 	response, err := mqcloudService.DeleteUser(deleteUserOptions)
		// 	if err != nil {
		// 		panic(err)
		// 	}
		// 	if response.StatusCode != 204 {
		// 		fmt.Printf("\nUnexpected response status code received from DeleteUser(): %d\n", response.StatusCode)
		// 	}

		// 	// end-delete_user

		// 	Expect(err).To(BeNil())
		// 	Expect(response.StatusCode).To(Equal(204))
		// })
		// It(`DeleteApplication request example`, func() {
		// 	// begin-delete_application
		// 	fmt.Println("\nDeleteApplication() result:")
		// 	deleteApplicationOptions := mqcloudService.NewDeleteApplicationOptions(
		// 		serviceinstance_guid,
		// 		application_id,
		// 	)

		// 	response, err := mqcloudService.DeleteApplication(deleteApplicationOptions)
		// 	if err != nil {
		// 		panic(err)
		// 	}
		// 	if response.StatusCode != 204 {
		// 		fmt.Printf("\nUnexpected response status code received from DeleteApplication(): %d\n", response.StatusCode)
		// 	}

		// 	// end-delete_application

		// 	Expect(err).To(BeNil())
		// 	Expect(response.StatusCode).To(Equal(204))
		// })

		// It(`DeleteKeyStoreCertificate request example`, func() {
		// 	// begin-delete_key_store_certificate
		// 	fmt.Println("\nDeleteKeyStoreCertificate() result:")
		// 	deleteKeyStoreCertificateOptions := mqcloudService.NewDeleteKeyStoreCertificateOptions(
		// 		serviceinstance_guid,
		// 		queue_manager_id,
		// 		keystore_certificate_id,
		// 	)

		// 	response, err := mqcloudService.DeleteKeyStoreCertificate(deleteKeyStoreCertificateOptions)
		// 	if err != nil {
		// 		panic(err)
		// 	}
		// 	if response.StatusCode != 204 {
		// 		fmt.Printf("\nUnexpected response status code received from DeleteKeyStoreCertificate(): %d\n", response.StatusCode)
		// 	}

		// 	// end-delete_key_store_certificate

		// 	Expect(err).To(BeNil())
		// 	Expect(response.StatusCode).To(Equal(204))
		// })
		It(`DeleteQueueManager request example`, func() {
			fmt.Println("\nDeleteQueueManager() result:")
			// begin-delete_queue_manager

			deleteQueueManagerOptions := mqcloudService.NewDeleteQueueManagerOptions(
				serviceinstance_guid,
				queue_manager_id,
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
