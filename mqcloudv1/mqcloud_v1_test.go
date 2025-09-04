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
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"time"

	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/mqcloud-go-sdk/mqcloudv1"
	"github.com/go-openapi/strfmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe(`MqcloudV1`, func() {
	var testServer *httptest.Server
	Describe(`Service constructor tests`, func() {
		It(`Instantiate service client`, func() {
			mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1(&mqcloudv1.MqcloudV1Options{
				Authenticator: &core.NoAuthAuthenticator{},
			})
			Expect(mqcloudService).ToNot(BeNil())
			Expect(serviceErr).To(BeNil())
		})
		It(`Instantiate service client with error: Invalid URL`, func() {
			mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1(&mqcloudv1.MqcloudV1Options{
				URL: "{BAD_URL_STRING",
			})
			Expect(mqcloudService).To(BeNil())
			Expect(serviceErr).ToNot(BeNil())
		})
		It(`Instantiate service client with error: Invalid Auth`, func() {
			mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1(&mqcloudv1.MqcloudV1Options{
				URL: "https://mqcloudv1/api",
				Authenticator: &core.BasicAuthenticator{
					Username: "",
					Password: "",
				},
			})
			Expect(mqcloudService).To(BeNil())
			Expect(serviceErr).ToNot(BeNil())
		})
	})
	Describe(`Service constructor tests using external config`, func() {
		Context(`Using external config, construct service client instances`, func() {
			// Map containing environment variables used in testing.
			var testEnvironment = map[string]string{
				"MQCLOUD_URL":       "https://mqcloudv1/api",
				"MQCLOUD_AUTH_TYPE": "noauth",
			}

			It(`Create service client using external config successfully`, func() {
				SetTestEnvironment(testEnvironment)
				mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1UsingExternalConfig(&mqcloudv1.MqcloudV1Options{})
				Expect(mqcloudService).ToNot(BeNil())
				Expect(serviceErr).To(BeNil())
				ClearTestEnvironment(testEnvironment)

				clone := mqcloudService.Clone()
				Expect(clone).ToNot(BeNil())
				Expect(clone.Service != mqcloudService.Service).To(BeTrue())
				Expect(clone.GetServiceURL()).To(Equal(mqcloudService.GetServiceURL()))
				Expect(clone.Service.Options.Authenticator).To(Equal(mqcloudService.Service.Options.Authenticator))
			})
			It(`Create service client using external config and set url from constructor successfully`, func() {
				SetTestEnvironment(testEnvironment)
				mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1UsingExternalConfig(&mqcloudv1.MqcloudV1Options{
					URL: "https://testService/api",
				})
				Expect(mqcloudService).ToNot(BeNil())
				Expect(serviceErr).To(BeNil())
				Expect(mqcloudService.Service.GetServiceURL()).To(Equal("https://testService/api"))
				ClearTestEnvironment(testEnvironment)

				clone := mqcloudService.Clone()
				Expect(clone).ToNot(BeNil())
				Expect(clone.Service != mqcloudService.Service).To(BeTrue())
				Expect(clone.GetServiceURL()).To(Equal(mqcloudService.GetServiceURL()))
				Expect(clone.Service.Options.Authenticator).To(Equal(mqcloudService.Service.Options.Authenticator))
			})
			It(`Create service client using external config and set url programatically successfully`, func() {
				SetTestEnvironment(testEnvironment)
				mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1UsingExternalConfig(&mqcloudv1.MqcloudV1Options{})
				err := mqcloudService.SetServiceURL("https://testService/api")
				Expect(err).To(BeNil())
				Expect(mqcloudService).ToNot(BeNil())
				Expect(serviceErr).To(BeNil())
				Expect(mqcloudService.Service.GetServiceURL()).To(Equal("https://testService/api"))
				ClearTestEnvironment(testEnvironment)

				clone := mqcloudService.Clone()
				Expect(clone).ToNot(BeNil())
				Expect(clone.Service != mqcloudService.Service).To(BeTrue())
				Expect(clone.GetServiceURL()).To(Equal(mqcloudService.GetServiceURL()))
				Expect(clone.Service.Options.Authenticator).To(Equal(mqcloudService.Service.Options.Authenticator))
			})
		})
		Context(`Using external config, construct service client instances with error: Invalid Auth`, func() {
			// Map containing environment variables used in testing.
			var testEnvironment = map[string]string{
				"MQCLOUD_URL":       "https://mqcloudv1/api",
				"MQCLOUD_AUTH_TYPE": "someOtherAuth",
			}

			SetTestEnvironment(testEnvironment)
			mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1UsingExternalConfig(&mqcloudv1.MqcloudV1Options{})

			It(`Instantiate service client with error`, func() {
				Expect(mqcloudService).To(BeNil())
				Expect(serviceErr).ToNot(BeNil())
				ClearTestEnvironment(testEnvironment)
			})
		})
		Context(`Using external config, construct service client instances with error: Invalid URL`, func() {
			// Map containing environment variables used in testing.
			var testEnvironment = map[string]string{
				"MQCLOUD_AUTH_TYPE": "NOAuth",
			}

			SetTestEnvironment(testEnvironment)
			mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1UsingExternalConfig(&mqcloudv1.MqcloudV1Options{
				URL: "{BAD_URL_STRING",
			})

			It(`Instantiate service client with error`, func() {
				Expect(mqcloudService).To(BeNil())
				Expect(serviceErr).ToNot(BeNil())
				ClearTestEnvironment(testEnvironment)
			})
		})
	})
	Describe(`Regional endpoint tests`, func() {
		It(`GetServiceURLForRegion(region string)`, func() {
			var url string
			var err error
			url, err = mqcloudv1.GetServiceURLForRegion("INVALID_REGION")
			Expect(url).To(BeEmpty())
			Expect(err).ToNot(BeNil())
			fmt.Fprintf(GinkgoWriter, "Expected error: %s\n", err.Error())
		})
	})
	Describe(`Parameterized URL tests`, func() {
		It(`Format parameterized URL with all default values`, func() {
			constructedURL, err := mqcloudv1.ConstructServiceURL(nil)
			Expect(constructedURL).To(Equal("https://api.private.eu-de.mq2.cloud.ibm.com"))
			Expect(constructedURL).ToNot(BeNil())
			Expect(err).To(BeNil())
		})
		It(`Return an error if a provided variable name is invalid`, func() {
			var providedUrlVariables = map[string]string{
				"invalid_variable_name": "value",
			}
			constructedURL, err := mqcloudv1.ConstructServiceURL(providedUrlVariables)
			Expect(constructedURL).To(Equal(""))
			Expect(err).ToNot(BeNil())
		})
	})
	Describe(`GetUsageDetails(getUsageDetailsOptions *GetUsageDetailsOptions) - Operation response error`, func() {
		acceptLanguage := "en-US,en;q=0.5"
		getUsageDetailsPath := "/v1/a2b4d4bc-dadb-4637-bcec-9b7d1e723af8/usage"
		Context(`Using mock server endpoint with invalid JSON response`, func() {
			BeforeEach(func() {
				testServer = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
					defer GinkgoRecover()

					// Verify the contents of the request
					Expect(req.URL.EscapedPath()).To(Equal(getUsageDetailsPath))
					Expect(req.Method).To(Equal("GET"))
					Expect(req.Header["Accept-Language"]).ToNot(BeNil())
					Expect(req.Header["Accept-Language"][0]).To(Equal(fmt.Sprintf("%v", "en-US,en;q=0.5")))
					res.Header().Set("Content-type", "application/json")
					res.WriteHeader(200)
					fmt.Fprint(res, `} this is not valid json {`)
				}))
			})
			It(`Invoke GetUsageDetails with error: Operation response processing error`, func() {
				mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1(&mqcloudv1.MqcloudV1Options{
					URL:            testServer.URL,
					Authenticator:  &core.NoAuthAuthenticator{},
					AcceptLanguage: core.StringPtr(acceptLanguage),
				})
				Expect(serviceErr).To(BeNil())
				Expect(mqcloudService).ToNot(BeNil())

				// Construct an instance of the GetUsageDetailsOptions model
				getUsageDetailsOptionsModel := new(mqcloudv1.GetUsageDetailsOptions)
				getUsageDetailsOptionsModel.ServiceInstanceGuid = core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				getUsageDetailsOptionsModel.Headers = map[string]string{"x-custom-header": "x-custom-value"}
				// Expect response parsing to fail since we are receiving a text/plain response
				result, response, operationErr := mqcloudService.GetUsageDetails(getUsageDetailsOptionsModel)
				Expect(operationErr).ToNot(BeNil())
				Expect(response).ToNot(BeNil())
				Expect(result).To(BeNil())

				// Enable retries and test again
				mqcloudService.EnableRetries(0, 0)
				result, response, operationErr = mqcloudService.GetUsageDetails(getUsageDetailsOptionsModel)
				Expect(operationErr).ToNot(BeNil())
				Expect(response).ToNot(BeNil())
				Expect(result).To(BeNil())
			})
			AfterEach(func() {
				testServer.Close()
			})
		})
	})
	Describe(`GetUsageDetails(getUsageDetailsOptions *GetUsageDetailsOptions)`, func() {
		acceptLanguage := "en-US,en;q=0.5"
		getUsageDetailsPath := "/v1/a2b4d4bc-dadb-4637-bcec-9b7d1e723af8/usage"
		Context(`Using mock server endpoint with timeout`, func() {
			BeforeEach(func() {
				testServer = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
					defer GinkgoRecover()

					// Verify the contents of the request
					Expect(req.URL.EscapedPath()).To(Equal(getUsageDetailsPath))
					Expect(req.Method).To(Equal("GET"))

					Expect(req.Header["Accept-Language"]).ToNot(BeNil())
					Expect(req.Header["Accept-Language"][0]).To(Equal(fmt.Sprintf("%v", "en-US,en;q=0.5")))
					// Sleep a short time to support a timeout test
					time.Sleep(100 * time.Millisecond)

					// Set mock response
					res.Header().Set("Content-type", "application/json")
					res.WriteHeader(200)
					fmt.Fprintf(res, "%s", `{"vpc_entitlement": 4.0, "vpc_usage": 3.3}`)
				}))
			})
			It(`Invoke GetUsageDetails successfully with retries`, func() {
				mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1(&mqcloudv1.MqcloudV1Options{
					URL:            testServer.URL,
					Authenticator:  &core.NoAuthAuthenticator{},
					AcceptLanguage: core.StringPtr(acceptLanguage),
				})
				Expect(serviceErr).To(BeNil())
				Expect(mqcloudService).ToNot(BeNil())
				mqcloudService.EnableRetries(0, 0)

				// Construct an instance of the GetUsageDetailsOptions model
				getUsageDetailsOptionsModel := new(mqcloudv1.GetUsageDetailsOptions)
				getUsageDetailsOptionsModel.ServiceInstanceGuid = core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				getUsageDetailsOptionsModel.Headers = map[string]string{"x-custom-header": "x-custom-value"}

				// Invoke operation with a Context to test a timeout error
				ctx, cancelFunc := context.WithTimeout(context.Background(), 80*time.Millisecond)
				defer cancelFunc()
				_, _, operationErr := mqcloudService.GetUsageDetailsWithContext(ctx, getUsageDetailsOptionsModel)
				Expect(operationErr).ToNot(BeNil())
				Expect(operationErr.Error()).To(ContainSubstring("deadline exceeded"))

				// Disable retries and test again
				mqcloudService.DisableRetries()
				result, response, operationErr := mqcloudService.GetUsageDetails(getUsageDetailsOptionsModel)
				Expect(operationErr).To(BeNil())
				Expect(response).ToNot(BeNil())
				Expect(result).ToNot(BeNil())

				// Re-test the timeout error with retries disabled
				ctx, cancelFunc2 := context.WithTimeout(context.Background(), 80*time.Millisecond)
				defer cancelFunc2()
				_, _, operationErr = mqcloudService.GetUsageDetailsWithContext(ctx, getUsageDetailsOptionsModel)
				Expect(operationErr).ToNot(BeNil())
				Expect(operationErr.Error()).To(ContainSubstring("deadline exceeded"))
			})
			AfterEach(func() {
				testServer.Close()
			})
		})
		Context(`Using mock server endpoint`, func() {
			BeforeEach(func() {
				testServer = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
					defer GinkgoRecover()

					// Verify the contents of the request
					Expect(req.URL.EscapedPath()).To(Equal(getUsageDetailsPath))
					Expect(req.Method).To(Equal("GET"))

					Expect(req.Header["Accept-Language"]).ToNot(BeNil())
					Expect(req.Header["Accept-Language"][0]).To(Equal(fmt.Sprintf("%v", "en-US,en;q=0.5")))
					// Set mock response
					res.Header().Set("Content-type", "application/json")
					res.WriteHeader(200)
					fmt.Fprintf(res, "%s", `{"vpc_entitlement": 4.0, "vpc_usage": 3.3}`)
				}))
			})
			It(`Invoke GetUsageDetails successfully`, func() {
				mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1(&mqcloudv1.MqcloudV1Options{
					URL:            testServer.URL,
					Authenticator:  &core.NoAuthAuthenticator{},
					AcceptLanguage: core.StringPtr(acceptLanguage),
				})
				Expect(serviceErr).To(BeNil())
				Expect(mqcloudService).ToNot(BeNil())

				// Invoke operation with nil options model (negative test)
				result, response, operationErr := mqcloudService.GetUsageDetails(nil)
				Expect(operationErr).NotTo(BeNil())
				Expect(response).To(BeNil())
				Expect(result).To(BeNil())

				// Construct an instance of the GetUsageDetailsOptions model
				getUsageDetailsOptionsModel := new(mqcloudv1.GetUsageDetailsOptions)
				getUsageDetailsOptionsModel.ServiceInstanceGuid = core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				getUsageDetailsOptionsModel.Headers = map[string]string{"x-custom-header": "x-custom-value"}

				// Invoke operation with valid options model (positive test)
				result, response, operationErr = mqcloudService.GetUsageDetails(getUsageDetailsOptionsModel)
				Expect(operationErr).To(BeNil())
				Expect(response).ToNot(BeNil())
				Expect(result).ToNot(BeNil())

			})
			It(`Invoke GetUsageDetails with error: Operation validation and request error`, func() {
				mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1(&mqcloudv1.MqcloudV1Options{
					URL:            testServer.URL,
					Authenticator:  &core.NoAuthAuthenticator{},
					AcceptLanguage: core.StringPtr(acceptLanguage),
				})
				Expect(serviceErr).To(BeNil())
				Expect(mqcloudService).ToNot(BeNil())

				// Construct an instance of the GetUsageDetailsOptions model
				getUsageDetailsOptionsModel := new(mqcloudv1.GetUsageDetailsOptions)
				getUsageDetailsOptionsModel.ServiceInstanceGuid = core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				getUsageDetailsOptionsModel.Headers = map[string]string{"x-custom-header": "x-custom-value"}
				// Invoke operation with empty URL (negative test)
				err := mqcloudService.SetServiceURL("")
				Expect(err).To(BeNil())
				result, response, operationErr := mqcloudService.GetUsageDetails(getUsageDetailsOptionsModel)
				Expect(operationErr).ToNot(BeNil())
				Expect(operationErr.Error()).To(ContainSubstring(core.ERRORMSG_SERVICE_URL_MISSING))
				Expect(response).To(BeNil())
				Expect(result).To(BeNil())
				// Construct a second instance of the GetUsageDetailsOptions model with no property values
				getUsageDetailsOptionsModelNew := new(mqcloudv1.GetUsageDetailsOptions)
				// Invoke operation with invalid model (negative test)
				result, response, operationErr = mqcloudService.GetUsageDetails(getUsageDetailsOptionsModelNew)
				Expect(operationErr).ToNot(BeNil())
				Expect(response).To(BeNil())
				Expect(result).To(BeNil())
			})
			AfterEach(func() {
				testServer.Close()
			})
		})
		Context(`Using mock server endpoint with missing response body`, func() {
			BeforeEach(func() {
				testServer = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
					defer GinkgoRecover()

					// Set success status code with no respoonse body
					res.WriteHeader(200)
				}))
			})
			It(`Invoke GetUsageDetails successfully`, func() {
				mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1(&mqcloudv1.MqcloudV1Options{
					URL:            testServer.URL,
					Authenticator:  &core.NoAuthAuthenticator{},
					AcceptLanguage: core.StringPtr(acceptLanguage),
				})
				Expect(serviceErr).To(BeNil())
				Expect(mqcloudService).ToNot(BeNil())

				// Construct an instance of the GetUsageDetailsOptions model
				getUsageDetailsOptionsModel := new(mqcloudv1.GetUsageDetailsOptions)
				getUsageDetailsOptionsModel.ServiceInstanceGuid = core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				getUsageDetailsOptionsModel.Headers = map[string]string{"x-custom-header": "x-custom-value"}

				// Invoke operation
				result, response, operationErr := mqcloudService.GetUsageDetails(getUsageDetailsOptionsModel)
				Expect(operationErr).To(BeNil())
				Expect(response).ToNot(BeNil())

				// Verify a nil result
				Expect(result).To(BeNil())
			})
			AfterEach(func() {
				testServer.Close()
			})
		})
	})
	Describe(`GetOptions(getOptionsOptions *GetOptionsOptions) - Operation response error`, func() {
		acceptLanguage := "en-US,en;q=0.5"
		getOptionsPath := "/v1/a2b4d4bc-dadb-4637-bcec-9b7d1e723af8/options"
		Context(`Using mock server endpoint with invalid JSON response`, func() {
			BeforeEach(func() {
				testServer = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
					defer GinkgoRecover()

					// Verify the contents of the request
					Expect(req.URL.EscapedPath()).To(Equal(getOptionsPath))
					Expect(req.Method).To(Equal("GET"))
					Expect(req.Header["Accept-Language"]).ToNot(BeNil())
					Expect(req.Header["Accept-Language"][0]).To(Equal(fmt.Sprintf("%v", "en-US,en;q=0.5")))
					res.Header().Set("Content-type", "application/json")
					res.WriteHeader(200)
					fmt.Fprint(res, `} this is not valid json {`)
				}))
			})
			It(`Invoke GetOptions with error: Operation response processing error`, func() {
				mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1(&mqcloudv1.MqcloudV1Options{
					URL:            testServer.URL,
					Authenticator:  &core.NoAuthAuthenticator{},
					AcceptLanguage: core.StringPtr(acceptLanguage),
				})
				Expect(serviceErr).To(BeNil())
				Expect(mqcloudService).ToNot(BeNil())

				// Construct an instance of the GetOptionsOptions model
				getOptionsOptionsModel := new(mqcloudv1.GetOptionsOptions)
				getOptionsOptionsModel.ServiceInstanceGuid = core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				getOptionsOptionsModel.Headers = map[string]string{"x-custom-header": "x-custom-value"}
				// Expect response parsing to fail since we are receiving a text/plain response
				result, response, operationErr := mqcloudService.GetOptions(getOptionsOptionsModel)
				Expect(operationErr).ToNot(BeNil())
				Expect(response).ToNot(BeNil())
				Expect(result).To(BeNil())

				// Enable retries and test again
				mqcloudService.EnableRetries(0, 0)
				result, response, operationErr = mqcloudService.GetOptions(getOptionsOptionsModel)
				Expect(operationErr).ToNot(BeNil())
				Expect(response).ToNot(BeNil())
				Expect(result).To(BeNil())
			})
			AfterEach(func() {
				testServer.Close()
			})
		})
	})
	Describe(`GetOptions(getOptionsOptions *GetOptionsOptions)`, func() {
		acceptLanguage := "en-US,en;q=0.5"
		getOptionsPath := "/v1/a2b4d4bc-dadb-4637-bcec-9b7d1e723af8/options"
		Context(`Using mock server endpoint with timeout`, func() {
			BeforeEach(func() {
				testServer = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
					defer GinkgoRecover()

					// Verify the contents of the request
					Expect(req.URL.EscapedPath()).To(Equal(getOptionsPath))
					Expect(req.Method).To(Equal("GET"))

					Expect(req.Header["Accept-Language"]).ToNot(BeNil())
					Expect(req.Header["Accept-Language"][0]).To(Equal(fmt.Sprintf("%v", "en-US,en;q=0.5")))
					// Sleep a short time to support a timeout test
					time.Sleep(100 * time.Millisecond)

					// Set mock response
					res.Header().Set("Content-type", "application/json")
					res.WriteHeader(200)
					fmt.Fprintf(res, "%s", `{"locations": ["reserved-eu-de-cluster-f884"], "sizes": ["small"], "versions": ["Versions"], "latest_version": "9.3.2_2"}`)
				}))
			})
			It(`Invoke GetOptions successfully with retries`, func() {
				mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1(&mqcloudv1.MqcloudV1Options{
					URL:            testServer.URL,
					Authenticator:  &core.NoAuthAuthenticator{},
					AcceptLanguage: core.StringPtr(acceptLanguage),
				})
				Expect(serviceErr).To(BeNil())
				Expect(mqcloudService).ToNot(BeNil())
				mqcloudService.EnableRetries(0, 0)

				// Construct an instance of the GetOptionsOptions model
				getOptionsOptionsModel := new(mqcloudv1.GetOptionsOptions)
				getOptionsOptionsModel.ServiceInstanceGuid = core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				getOptionsOptionsModel.Headers = map[string]string{"x-custom-header": "x-custom-value"}

				// Invoke operation with a Context to test a timeout error
				ctx, cancelFunc := context.WithTimeout(context.Background(), 80*time.Millisecond)
				defer cancelFunc()
				_, _, operationErr := mqcloudService.GetOptionsWithContext(ctx, getOptionsOptionsModel)
				Expect(operationErr).ToNot(BeNil())
				Expect(operationErr.Error()).To(ContainSubstring("deadline exceeded"))

				// Disable retries and test again
				mqcloudService.DisableRetries()
				result, response, operationErr := mqcloudService.GetOptions(getOptionsOptionsModel)
				Expect(operationErr).To(BeNil())
				Expect(response).ToNot(BeNil())
				Expect(result).ToNot(BeNil())

				// Re-test the timeout error with retries disabled
				ctx, cancelFunc2 := context.WithTimeout(context.Background(), 80*time.Millisecond)
				defer cancelFunc2()
				_, _, operationErr = mqcloudService.GetOptionsWithContext(ctx, getOptionsOptionsModel)
				Expect(operationErr).ToNot(BeNil())
				Expect(operationErr.Error()).To(ContainSubstring("deadline exceeded"))
			})
			AfterEach(func() {
				testServer.Close()
			})
		})
		Context(`Using mock server endpoint`, func() {
			BeforeEach(func() {
				testServer = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
					defer GinkgoRecover()

					// Verify the contents of the request
					Expect(req.URL.EscapedPath()).To(Equal(getOptionsPath))
					Expect(req.Method).To(Equal("GET"))

					Expect(req.Header["Accept-Language"]).ToNot(BeNil())
					Expect(req.Header["Accept-Language"][0]).To(Equal(fmt.Sprintf("%v", "en-US,en;q=0.5")))
					// Set mock response
					res.Header().Set("Content-type", "application/json")
					res.WriteHeader(200)
					fmt.Fprintf(res, "%s", `{"locations": ["reserved-eu-de-cluster-f884"], "sizes": ["small"], "versions": ["Versions"], "latest_version": "9.3.2_2"}`)
				}))
			})
			It(`Invoke GetOptions successfully`, func() {
				mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1(&mqcloudv1.MqcloudV1Options{
					URL:            testServer.URL,
					Authenticator:  &core.NoAuthAuthenticator{},
					AcceptLanguage: core.StringPtr(acceptLanguage),
				})
				Expect(serviceErr).To(BeNil())
				Expect(mqcloudService).ToNot(BeNil())

				// Invoke operation with nil options model (negative test)
				result, response, operationErr := mqcloudService.GetOptions(nil)
				Expect(operationErr).NotTo(BeNil())
				Expect(response).To(BeNil())
				Expect(result).To(BeNil())

				// Construct an instance of the GetOptionsOptions model
				getOptionsOptionsModel := new(mqcloudv1.GetOptionsOptions)
				getOptionsOptionsModel.ServiceInstanceGuid = core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				getOptionsOptionsModel.Headers = map[string]string{"x-custom-header": "x-custom-value"}

				// Invoke operation with valid options model (positive test)
				result, response, operationErr = mqcloudService.GetOptions(getOptionsOptionsModel)
				Expect(operationErr).To(BeNil())
				Expect(response).ToNot(BeNil())
				Expect(result).ToNot(BeNil())

			})
			It(`Invoke GetOptions with error: Operation validation and request error`, func() {
				mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1(&mqcloudv1.MqcloudV1Options{
					URL:            testServer.URL,
					Authenticator:  &core.NoAuthAuthenticator{},
					AcceptLanguage: core.StringPtr(acceptLanguage),
				})
				Expect(serviceErr).To(BeNil())
				Expect(mqcloudService).ToNot(BeNil())

				// Construct an instance of the GetOptionsOptions model
				getOptionsOptionsModel := new(mqcloudv1.GetOptionsOptions)
				getOptionsOptionsModel.ServiceInstanceGuid = core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				getOptionsOptionsModel.Headers = map[string]string{"x-custom-header": "x-custom-value"}
				// Invoke operation with empty URL (negative test)
				err := mqcloudService.SetServiceURL("")
				Expect(err).To(BeNil())
				result, response, operationErr := mqcloudService.GetOptions(getOptionsOptionsModel)
				Expect(operationErr).ToNot(BeNil())
				Expect(operationErr.Error()).To(ContainSubstring(core.ERRORMSG_SERVICE_URL_MISSING))
				Expect(response).To(BeNil())
				Expect(result).To(BeNil())
				// Construct a second instance of the GetOptionsOptions model with no property values
				getOptionsOptionsModelNew := new(mqcloudv1.GetOptionsOptions)
				// Invoke operation with invalid model (negative test)
				result, response, operationErr = mqcloudService.GetOptions(getOptionsOptionsModelNew)
				Expect(operationErr).ToNot(BeNil())
				Expect(response).To(BeNil())
				Expect(result).To(BeNil())
			})
			AfterEach(func() {
				testServer.Close()
			})
		})
		Context(`Using mock server endpoint with missing response body`, func() {
			BeforeEach(func() {
				testServer = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
					defer GinkgoRecover()

					// Set success status code with no respoonse body
					res.WriteHeader(200)
				}))
			})
			It(`Invoke GetOptions successfully`, func() {
				mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1(&mqcloudv1.MqcloudV1Options{
					URL:            testServer.URL,
					Authenticator:  &core.NoAuthAuthenticator{},
					AcceptLanguage: core.StringPtr(acceptLanguage),
				})
				Expect(serviceErr).To(BeNil())
				Expect(mqcloudService).ToNot(BeNil())

				// Construct an instance of the GetOptionsOptions model
				getOptionsOptionsModel := new(mqcloudv1.GetOptionsOptions)
				getOptionsOptionsModel.ServiceInstanceGuid = core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				getOptionsOptionsModel.Headers = map[string]string{"x-custom-header": "x-custom-value"}

				// Invoke operation
				result, response, operationErr := mqcloudService.GetOptions(getOptionsOptionsModel)
				Expect(operationErr).To(BeNil())
				Expect(response).ToNot(BeNil())

				// Verify a nil result
				Expect(result).To(BeNil())
			})
			AfterEach(func() {
				testServer.Close()
			})
		})
	})
	Describe(`CreateQueueManager(createQueueManagerOptions *CreateQueueManagerOptions) - Operation response error`, func() {
		acceptLanguage := "en-US,en;q=0.5"
		createQueueManagerPath := "/v1/a2b4d4bc-dadb-4637-bcec-9b7d1e723af8/queue_managers"
		Context(`Using mock server endpoint with invalid JSON response`, func() {
			BeforeEach(func() {
				testServer = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
					defer GinkgoRecover()

					// Verify the contents of the request
					Expect(req.URL.EscapedPath()).To(Equal(createQueueManagerPath))
					Expect(req.Method).To(Equal("POST"))
					Expect(req.Header["Accept-Language"]).ToNot(BeNil())
					Expect(req.Header["Accept-Language"][0]).To(Equal(fmt.Sprintf("%v", "en-US,en;q=0.5")))
					res.Header().Set("Content-type", "application/json")
					res.WriteHeader(202)
					fmt.Fprint(res, `} this is not valid json {`)
				}))
			})
			It(`Invoke CreateQueueManager with error: Operation response processing error`, func() {
				mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1(&mqcloudv1.MqcloudV1Options{
					URL:            testServer.URL,
					Authenticator:  &core.NoAuthAuthenticator{},
					AcceptLanguage: core.StringPtr(acceptLanguage),
				})
				Expect(serviceErr).To(BeNil())
				Expect(mqcloudService).ToNot(BeNil())

				// Construct an instance of the CreateQueueManagerOptions model
				createQueueManagerOptionsModel := new(mqcloudv1.CreateQueueManagerOptions)
				createQueueManagerOptionsModel.ServiceInstanceGuid = core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				createQueueManagerOptionsModel.Name = core.StringPtr("testqm")
				createQueueManagerOptionsModel.Location = core.StringPtr("reserved-eu-de-cluster-f884")
				createQueueManagerOptionsModel.Size = core.StringPtr("small")
				createQueueManagerOptionsModel.DisplayName = core.StringPtr("A test queue manager")
				createQueueManagerOptionsModel.Version = core.StringPtr("9.3.2_2")
				createQueueManagerOptionsModel.Headers = map[string]string{"x-custom-header": "x-custom-value"}
				// Expect response parsing to fail since we are receiving a text/plain response
				result, response, operationErr := mqcloudService.CreateQueueManager(createQueueManagerOptionsModel)
				Expect(operationErr).ToNot(BeNil())
				Expect(response).ToNot(BeNil())
				Expect(result).To(BeNil())

				// Enable retries and test again
				mqcloudService.EnableRetries(0, 0)
				result, response, operationErr = mqcloudService.CreateQueueManager(createQueueManagerOptionsModel)
				Expect(operationErr).ToNot(BeNil())
				Expect(response).ToNot(BeNil())
				Expect(result).To(BeNil())
			})
			AfterEach(func() {
				testServer.Close()
			})
		})
	})
	Describe(`CreateQueueManager(createQueueManagerOptions *CreateQueueManagerOptions)`, func() {
		acceptLanguage := "en-US,en;q=0.5"
		createQueueManagerPath := "/v1/a2b4d4bc-dadb-4637-bcec-9b7d1e723af8/queue_managers"
		Context(`Using mock server endpoint with timeout`, func() {
			BeforeEach(func() {
				testServer = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
					defer GinkgoRecover()

					// Verify the contents of the request
					Expect(req.URL.EscapedPath()).To(Equal(createQueueManagerPath))
					Expect(req.Method).To(Equal("POST"))

					// For gzip-disabled operation, verify Content-Encoding is not set.
					Expect(req.Header.Get("Content-Encoding")).To(BeEmpty())

					// If there is a body, then make sure we can read it
					bodyBuf := new(bytes.Buffer)
					if req.Header.Get("Content-Encoding") == "gzip" {
						body, err := core.NewGzipDecompressionReader(req.Body)
						Expect(err).To(BeNil())
						_, err = bodyBuf.ReadFrom(body)
						Expect(err).To(BeNil())
					} else {
						_, err := bodyBuf.ReadFrom(req.Body)
						Expect(err).To(BeNil())
					}
					fmt.Fprintf(GinkgoWriter, "  Request body: %s", bodyBuf.String())

					Expect(req.Header["Accept-Language"]).ToNot(BeNil())
					Expect(req.Header["Accept-Language"][0]).To(Equal(fmt.Sprintf("%v", "en-US,en;q=0.5")))
					// Sleep a short time to support a timeout test
					time.Sleep(100 * time.Millisecond)

					// Set mock response
					res.Header().Set("Content-type", "application/json")
					res.WriteHeader(202)
					fmt.Fprintf(res, "%s", `{"queue_manager_uri": "QueueManagerURI", "queue_manager_status_uri": "QueueManagerStatusURI", "queue_manager_id": "QueueManagerID"}`)
				}))
			})
			It(`Invoke CreateQueueManager successfully with retries`, func() {
				mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1(&mqcloudv1.MqcloudV1Options{
					URL:            testServer.URL,
					Authenticator:  &core.NoAuthAuthenticator{},
					AcceptLanguage: core.StringPtr(acceptLanguage),
				})
				Expect(serviceErr).To(BeNil())
				Expect(mqcloudService).ToNot(BeNil())
				mqcloudService.EnableRetries(0, 0)

				// Construct an instance of the CreateQueueManagerOptions model
				createQueueManagerOptionsModel := new(mqcloudv1.CreateQueueManagerOptions)
				createQueueManagerOptionsModel.ServiceInstanceGuid = core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				createQueueManagerOptionsModel.Name = core.StringPtr("testqm")
				createQueueManagerOptionsModel.Location = core.StringPtr("reserved-eu-de-cluster-f884")
				createQueueManagerOptionsModel.Size = core.StringPtr("small")
				createQueueManagerOptionsModel.DisplayName = core.StringPtr("A test queue manager")
				createQueueManagerOptionsModel.Version = core.StringPtr("9.3.2_2")
				createQueueManagerOptionsModel.Headers = map[string]string{"x-custom-header": "x-custom-value"}

				// Invoke operation with a Context to test a timeout error
				ctx, cancelFunc := context.WithTimeout(context.Background(), 80*time.Millisecond)
				defer cancelFunc()
				_, _, operationErr := mqcloudService.CreateQueueManagerWithContext(ctx, createQueueManagerOptionsModel)
				Expect(operationErr).ToNot(BeNil())
				Expect(operationErr.Error()).To(ContainSubstring("deadline exceeded"))

				// Disable retries and test again
				mqcloudService.DisableRetries()
				result, response, operationErr := mqcloudService.CreateQueueManager(createQueueManagerOptionsModel)
				Expect(operationErr).To(BeNil())
				Expect(response).ToNot(BeNil())
				Expect(result).ToNot(BeNil())

				// Re-test the timeout error with retries disabled
				ctx, cancelFunc2 := context.WithTimeout(context.Background(), 80*time.Millisecond)
				defer cancelFunc2()
				_, _, operationErr = mqcloudService.CreateQueueManagerWithContext(ctx, createQueueManagerOptionsModel)
				Expect(operationErr).ToNot(BeNil())
				Expect(operationErr.Error()).To(ContainSubstring("deadline exceeded"))
			})
			AfterEach(func() {
				testServer.Close()
			})
		})
		Context(`Using mock server endpoint`, func() {
			BeforeEach(func() {
				testServer = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
					defer GinkgoRecover()

					// Verify the contents of the request
					Expect(req.URL.EscapedPath()).To(Equal(createQueueManagerPath))
					Expect(req.Method).To(Equal("POST"))

					// For gzip-disabled operation, verify Content-Encoding is not set.
					Expect(req.Header.Get("Content-Encoding")).To(BeEmpty())

					// If there is a body, then make sure we can read it
					bodyBuf := new(bytes.Buffer)
					if req.Header.Get("Content-Encoding") == "gzip" {
						body, err := core.NewGzipDecompressionReader(req.Body)
						Expect(err).To(BeNil())
						_, err = bodyBuf.ReadFrom(body)
						Expect(err).To(BeNil())
					} else {
						_, err := bodyBuf.ReadFrom(req.Body)
						Expect(err).To(BeNil())
					}
					fmt.Fprintf(GinkgoWriter, "  Request body: %s", bodyBuf.String())

					Expect(req.Header["Accept-Language"]).ToNot(BeNil())
					Expect(req.Header["Accept-Language"][0]).To(Equal(fmt.Sprintf("%v", "en-US,en;q=0.5")))
					// Set mock response
					res.Header().Set("Content-type", "application/json")
					res.WriteHeader(202)
					fmt.Fprintf(res, "%s", `{"queue_manager_uri": "QueueManagerURI", "queue_manager_status_uri": "QueueManagerStatusURI", "queue_manager_id": "QueueManagerID"}`)
				}))
			})
			It(`Invoke CreateQueueManager successfully`, func() {
				mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1(&mqcloudv1.MqcloudV1Options{
					URL:            testServer.URL,
					Authenticator:  &core.NoAuthAuthenticator{},
					AcceptLanguage: core.StringPtr(acceptLanguage),
				})
				Expect(serviceErr).To(BeNil())
				Expect(mqcloudService).ToNot(BeNil())

				// Invoke operation with nil options model (negative test)
				result, response, operationErr := mqcloudService.CreateQueueManager(nil)
				Expect(operationErr).NotTo(BeNil())
				Expect(response).To(BeNil())
				Expect(result).To(BeNil())

				// Construct an instance of the CreateQueueManagerOptions model
				createQueueManagerOptionsModel := new(mqcloudv1.CreateQueueManagerOptions)
				createQueueManagerOptionsModel.ServiceInstanceGuid = core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				createQueueManagerOptionsModel.Name = core.StringPtr("testqm")
				createQueueManagerOptionsModel.Location = core.StringPtr("reserved-eu-de-cluster-f884")
				createQueueManagerOptionsModel.Size = core.StringPtr("small")
				createQueueManagerOptionsModel.DisplayName = core.StringPtr("A test queue manager")
				createQueueManagerOptionsModel.Version = core.StringPtr("9.3.2_2")
				createQueueManagerOptionsModel.Headers = map[string]string{"x-custom-header": "x-custom-value"}

				// Invoke operation with valid options model (positive test)
				result, response, operationErr = mqcloudService.CreateQueueManager(createQueueManagerOptionsModel)
				Expect(operationErr).To(BeNil())
				Expect(response).ToNot(BeNil())
				Expect(result).ToNot(BeNil())

			})
			It(`Invoke CreateQueueManager with error: Operation validation and request error`, func() {
				mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1(&mqcloudv1.MqcloudV1Options{
					URL:            testServer.URL,
					Authenticator:  &core.NoAuthAuthenticator{},
					AcceptLanguage: core.StringPtr(acceptLanguage),
				})
				Expect(serviceErr).To(BeNil())
				Expect(mqcloudService).ToNot(BeNil())

				// Construct an instance of the CreateQueueManagerOptions model
				createQueueManagerOptionsModel := new(mqcloudv1.CreateQueueManagerOptions)
				createQueueManagerOptionsModel.ServiceInstanceGuid = core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				createQueueManagerOptionsModel.Name = core.StringPtr("testqm")
				createQueueManagerOptionsModel.Location = core.StringPtr("reserved-eu-de-cluster-f884")
				createQueueManagerOptionsModel.Size = core.StringPtr("small")
				createQueueManagerOptionsModel.DisplayName = core.StringPtr("A test queue manager")
				createQueueManagerOptionsModel.Version = core.StringPtr("9.3.2_2")
				createQueueManagerOptionsModel.Headers = map[string]string{"x-custom-header": "x-custom-value"}
				// Invoke operation with empty URL (negative test)
				err := mqcloudService.SetServiceURL("")
				Expect(err).To(BeNil())
				result, response, operationErr := mqcloudService.CreateQueueManager(createQueueManagerOptionsModel)
				Expect(operationErr).ToNot(BeNil())
				Expect(operationErr.Error()).To(ContainSubstring(core.ERRORMSG_SERVICE_URL_MISSING))
				Expect(response).To(BeNil())
				Expect(result).To(BeNil())
				// Construct a second instance of the CreateQueueManagerOptions model with no property values
				createQueueManagerOptionsModelNew := new(mqcloudv1.CreateQueueManagerOptions)
				// Invoke operation with invalid model (negative test)
				result, response, operationErr = mqcloudService.CreateQueueManager(createQueueManagerOptionsModelNew)
				Expect(operationErr).ToNot(BeNil())
				Expect(response).To(BeNil())
				Expect(result).To(BeNil())
			})
			AfterEach(func() {
				testServer.Close()
			})
		})
		Context(`Using mock server endpoint with missing response body`, func() {
			BeforeEach(func() {
				testServer = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
					defer GinkgoRecover()

					// Set success status code with no respoonse body
					res.WriteHeader(202)
				}))
			})
			It(`Invoke CreateQueueManager successfully`, func() {
				mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1(&mqcloudv1.MqcloudV1Options{
					URL:            testServer.URL,
					Authenticator:  &core.NoAuthAuthenticator{},
					AcceptLanguage: core.StringPtr(acceptLanguage),
				})
				Expect(serviceErr).To(BeNil())
				Expect(mqcloudService).ToNot(BeNil())

				// Construct an instance of the CreateQueueManagerOptions model
				createQueueManagerOptionsModel := new(mqcloudv1.CreateQueueManagerOptions)
				createQueueManagerOptionsModel.ServiceInstanceGuid = core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				createQueueManagerOptionsModel.Name = core.StringPtr("testqm")
				createQueueManagerOptionsModel.Location = core.StringPtr("reserved-eu-de-cluster-f884")
				createQueueManagerOptionsModel.Size = core.StringPtr("small")
				createQueueManagerOptionsModel.DisplayName = core.StringPtr("A test queue manager")
				createQueueManagerOptionsModel.Version = core.StringPtr("9.3.2_2")
				createQueueManagerOptionsModel.Headers = map[string]string{"x-custom-header": "x-custom-value"}

				// Invoke operation
				result, response, operationErr := mqcloudService.CreateQueueManager(createQueueManagerOptionsModel)
				Expect(operationErr).To(BeNil())
				Expect(response).ToNot(BeNil())

				// Verify a nil result
				Expect(result).To(BeNil())
			})
			AfterEach(func() {
				testServer.Close()
			})
		})
	})
	Describe(`ListQueueManagers(listQueueManagersOptions *ListQueueManagersOptions) - Operation response error`, func() {
		acceptLanguage := "en-US,en;q=0.5"
		listQueueManagersPath := "/v1/a2b4d4bc-dadb-4637-bcec-9b7d1e723af8/queue_managers"
		Context(`Using mock server endpoint with invalid JSON response`, func() {
			BeforeEach(func() {
				testServer = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
					defer GinkgoRecover()

					// Verify the contents of the request
					Expect(req.URL.EscapedPath()).To(Equal(listQueueManagersPath))
					Expect(req.Method).To(Equal("GET"))
					Expect(req.Header["Accept-Language"]).ToNot(BeNil())
					Expect(req.Header["Accept-Language"][0]).To(Equal(fmt.Sprintf("%v", "en-US,en;q=0.5")))
					Expect(req.URL.Query()["offset"]).To(Equal([]string{fmt.Sprint(int64(0))}))
					Expect(req.URL.Query()["limit"]).To(Equal([]string{fmt.Sprint(int64(10))}))
					res.Header().Set("Content-type", "application/json")
					res.WriteHeader(200)
					fmt.Fprint(res, `} this is not valid json {`)
				}))
			})
			It(`Invoke ListQueueManagers with error: Operation response processing error`, func() {
				mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1(&mqcloudv1.MqcloudV1Options{
					URL:            testServer.URL,
					Authenticator:  &core.NoAuthAuthenticator{},
					AcceptLanguage: core.StringPtr(acceptLanguage),
				})
				Expect(serviceErr).To(BeNil())
				Expect(mqcloudService).ToNot(BeNil())

				// Construct an instance of the ListQueueManagersOptions model
				listQueueManagersOptionsModel := new(mqcloudv1.ListQueueManagersOptions)
				listQueueManagersOptionsModel.ServiceInstanceGuid = core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				listQueueManagersOptionsModel.Offset = core.Int64Ptr(int64(0))
				listQueueManagersOptionsModel.Limit = core.Int64Ptr(int64(10))
				listQueueManagersOptionsModel.Headers = map[string]string{"x-custom-header": "x-custom-value"}
				// Expect response parsing to fail since we are receiving a text/plain response
				result, response, operationErr := mqcloudService.ListQueueManagers(listQueueManagersOptionsModel)
				Expect(operationErr).ToNot(BeNil())
				Expect(response).ToNot(BeNil())
				Expect(result).To(BeNil())

				// Enable retries and test again
				mqcloudService.EnableRetries(0, 0)
				result, response, operationErr = mqcloudService.ListQueueManagers(listQueueManagersOptionsModel)
				Expect(operationErr).ToNot(BeNil())
				Expect(response).ToNot(BeNil())
				Expect(result).To(BeNil())
			})
			AfterEach(func() {
				testServer.Close()
			})
		})
	})
	Describe(`ListQueueManagers(listQueueManagersOptions *ListQueueManagersOptions)`, func() {
		acceptLanguage := "en-US,en;q=0.5"
		listQueueManagersPath := "/v1/a2b4d4bc-dadb-4637-bcec-9b7d1e723af8/queue_managers"
		Context(`Using mock server endpoint with timeout`, func() {
			BeforeEach(func() {
				testServer = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
					defer GinkgoRecover()

					// Verify the contents of the request
					Expect(req.URL.EscapedPath()).To(Equal(listQueueManagersPath))
					Expect(req.Method).To(Equal("GET"))

					Expect(req.Header["Accept-Language"]).ToNot(BeNil())
					Expect(req.Header["Accept-Language"][0]).To(Equal(fmt.Sprintf("%v", "en-US,en;q=0.5")))
					Expect(req.URL.Query()["offset"]).To(Equal([]string{fmt.Sprint(int64(0))}))
					Expect(req.URL.Query()["limit"]).To(Equal([]string{fmt.Sprint(int64(10))}))
					// Sleep a short time to support a timeout test
					time.Sleep(100 * time.Millisecond)

					// Set mock response
					res.Header().Set("Content-type", "application/json")
					res.WriteHeader(200)
					fmt.Fprintf(res, "%s", `{"offset": 6, "limit": 50, "first": {"href": "Href"}, "next": {"href": "Href"}, "previous": {"href": "Href"}, "queue_managers": [{"id": "ID", "name": "Name", "display_name": "DisplayName", "location": "reserved-eu-de-cluster-f884", "size": "small", "status_uri": "StatusURI", "version": "9.3.2_2", "web_console_url": "WebConsoleURL", "rest_api_endpoint_url": "RestApiEndpointURL", "administrator_api_endpoint_url": "AdministratorApiEndpointURL", "connection_info_uri": "ConnectionInfoURI", "date_created": "2020-01-13T15:39:35.000Z", "upgrade_available": true, "available_upgrade_versions_uri": "AvailableUpgradeVersionsURI", "href": "Href"}]}`)
				}))
			})
			It(`Invoke ListQueueManagers successfully with retries`, func() {
				mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1(&mqcloudv1.MqcloudV1Options{
					URL:            testServer.URL,
					Authenticator:  &core.NoAuthAuthenticator{},
					AcceptLanguage: core.StringPtr(acceptLanguage),
				})
				Expect(serviceErr).To(BeNil())
				Expect(mqcloudService).ToNot(BeNil())
				mqcloudService.EnableRetries(0, 0)

				// Construct an instance of the ListQueueManagersOptions model
				listQueueManagersOptionsModel := new(mqcloudv1.ListQueueManagersOptions)
				listQueueManagersOptionsModel.ServiceInstanceGuid = core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				listQueueManagersOptionsModel.Offset = core.Int64Ptr(int64(0))
				listQueueManagersOptionsModel.Limit = core.Int64Ptr(int64(10))
				listQueueManagersOptionsModel.Headers = map[string]string{"x-custom-header": "x-custom-value"}

				// Invoke operation with a Context to test a timeout error
				ctx, cancelFunc := context.WithTimeout(context.Background(), 80*time.Millisecond)
				defer cancelFunc()
				_, _, operationErr := mqcloudService.ListQueueManagersWithContext(ctx, listQueueManagersOptionsModel)
				Expect(operationErr).ToNot(BeNil())
				Expect(operationErr.Error()).To(ContainSubstring("deadline exceeded"))

				// Disable retries and test again
				mqcloudService.DisableRetries()
				result, response, operationErr := mqcloudService.ListQueueManagers(listQueueManagersOptionsModel)
				Expect(operationErr).To(BeNil())
				Expect(response).ToNot(BeNil())
				Expect(result).ToNot(BeNil())

				// Re-test the timeout error with retries disabled
				ctx, cancelFunc2 := context.WithTimeout(context.Background(), 80*time.Millisecond)
				defer cancelFunc2()
				_, _, operationErr = mqcloudService.ListQueueManagersWithContext(ctx, listQueueManagersOptionsModel)
				Expect(operationErr).ToNot(BeNil())
				Expect(operationErr.Error()).To(ContainSubstring("deadline exceeded"))
			})
			AfterEach(func() {
				testServer.Close()
			})
		})
		Context(`Using mock server endpoint`, func() {
			BeforeEach(func() {
				testServer = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
					defer GinkgoRecover()

					// Verify the contents of the request
					Expect(req.URL.EscapedPath()).To(Equal(listQueueManagersPath))
					Expect(req.Method).To(Equal("GET"))

					Expect(req.Header["Accept-Language"]).ToNot(BeNil())
					Expect(req.Header["Accept-Language"][0]).To(Equal(fmt.Sprintf("%v", "en-US,en;q=0.5")))
					Expect(req.URL.Query()["offset"]).To(Equal([]string{fmt.Sprint(int64(0))}))
					Expect(req.URL.Query()["limit"]).To(Equal([]string{fmt.Sprint(int64(10))}))
					// Set mock response
					res.Header().Set("Content-type", "application/json")
					res.WriteHeader(200)
					fmt.Fprintf(res, "%s", `{"offset": 6, "limit": 50, "first": {"href": "Href"}, "next": {"href": "Href"}, "previous": {"href": "Href"}, "queue_managers": [{"id": "ID", "name": "Name", "display_name": "DisplayName", "location": "reserved-eu-de-cluster-f884", "size": "small", "status_uri": "StatusURI", "version": "9.3.2_2", "web_console_url": "WebConsoleURL", "rest_api_endpoint_url": "RestApiEndpointURL", "administrator_api_endpoint_url": "AdministratorApiEndpointURL", "connection_info_uri": "ConnectionInfoURI", "date_created": "2020-01-13T15:39:35.000Z", "upgrade_available": true, "available_upgrade_versions_uri": "AvailableUpgradeVersionsURI", "href": "Href"}]}`)
				}))
			})
			It(`Invoke ListQueueManagers successfully`, func() {
				mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1(&mqcloudv1.MqcloudV1Options{
					URL:            testServer.URL,
					Authenticator:  &core.NoAuthAuthenticator{},
					AcceptLanguage: core.StringPtr(acceptLanguage),
				})
				Expect(serviceErr).To(BeNil())
				Expect(mqcloudService).ToNot(BeNil())

				// Invoke operation with nil options model (negative test)
				result, response, operationErr := mqcloudService.ListQueueManagers(nil)
				Expect(operationErr).NotTo(BeNil())
				Expect(response).To(BeNil())
				Expect(result).To(BeNil())

				// Construct an instance of the ListQueueManagersOptions model
				listQueueManagersOptionsModel := new(mqcloudv1.ListQueueManagersOptions)
				listQueueManagersOptionsModel.ServiceInstanceGuid = core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				listQueueManagersOptionsModel.Offset = core.Int64Ptr(int64(0))
				listQueueManagersOptionsModel.Limit = core.Int64Ptr(int64(10))
				listQueueManagersOptionsModel.Headers = map[string]string{"x-custom-header": "x-custom-value"}

				// Invoke operation with valid options model (positive test)
				result, response, operationErr = mqcloudService.ListQueueManagers(listQueueManagersOptionsModel)
				Expect(operationErr).To(BeNil())
				Expect(response).ToNot(BeNil())
				Expect(result).ToNot(BeNil())

			})
			It(`Invoke ListQueueManagers with error: Operation validation and request error`, func() {
				mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1(&mqcloudv1.MqcloudV1Options{
					URL:            testServer.URL,
					Authenticator:  &core.NoAuthAuthenticator{},
					AcceptLanguage: core.StringPtr(acceptLanguage),
				})
				Expect(serviceErr).To(BeNil())
				Expect(mqcloudService).ToNot(BeNil())

				// Construct an instance of the ListQueueManagersOptions model
				listQueueManagersOptionsModel := new(mqcloudv1.ListQueueManagersOptions)
				listQueueManagersOptionsModel.ServiceInstanceGuid = core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				listQueueManagersOptionsModel.Offset = core.Int64Ptr(int64(0))
				listQueueManagersOptionsModel.Limit = core.Int64Ptr(int64(10))
				listQueueManagersOptionsModel.Headers = map[string]string{"x-custom-header": "x-custom-value"}
				// Invoke operation with empty URL (negative test)
				err := mqcloudService.SetServiceURL("")
				Expect(err).To(BeNil())
				result, response, operationErr := mqcloudService.ListQueueManagers(listQueueManagersOptionsModel)
				Expect(operationErr).ToNot(BeNil())
				Expect(operationErr.Error()).To(ContainSubstring(core.ERRORMSG_SERVICE_URL_MISSING))
				Expect(response).To(BeNil())
				Expect(result).To(BeNil())
				// Construct a second instance of the ListQueueManagersOptions model with no property values
				listQueueManagersOptionsModelNew := new(mqcloudv1.ListQueueManagersOptions)
				// Invoke operation with invalid model (negative test)
				result, response, operationErr = mqcloudService.ListQueueManagers(listQueueManagersOptionsModelNew)
				Expect(operationErr).ToNot(BeNil())
				Expect(response).To(BeNil())
				Expect(result).To(BeNil())
			})
			AfterEach(func() {
				testServer.Close()
			})
		})
		Context(`Using mock server endpoint with missing response body`, func() {
			BeforeEach(func() {
				testServer = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
					defer GinkgoRecover()

					// Set success status code with no respoonse body
					res.WriteHeader(200)
				}))
			})
			It(`Invoke ListQueueManagers successfully`, func() {
				mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1(&mqcloudv1.MqcloudV1Options{
					URL:            testServer.URL,
					Authenticator:  &core.NoAuthAuthenticator{},
					AcceptLanguage: core.StringPtr(acceptLanguage),
				})
				Expect(serviceErr).To(BeNil())
				Expect(mqcloudService).ToNot(BeNil())

				// Construct an instance of the ListQueueManagersOptions model
				listQueueManagersOptionsModel := new(mqcloudv1.ListQueueManagersOptions)
				listQueueManagersOptionsModel.ServiceInstanceGuid = core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				listQueueManagersOptionsModel.Offset = core.Int64Ptr(int64(0))
				listQueueManagersOptionsModel.Limit = core.Int64Ptr(int64(10))
				listQueueManagersOptionsModel.Headers = map[string]string{"x-custom-header": "x-custom-value"}

				// Invoke operation
				result, response, operationErr := mqcloudService.ListQueueManagers(listQueueManagersOptionsModel)
				Expect(operationErr).To(BeNil())
				Expect(response).ToNot(BeNil())

				// Verify a nil result
				Expect(result).To(BeNil())
			})
			AfterEach(func() {
				testServer.Close()
			})
		})
		Context(`Test pagination helper method on response`, func() {
			It(`Invoke GetNextOffset successfully`, func() {
				responseObject := new(mqcloudv1.QueueManagerDetailsCollection)
				nextObject := new(mqcloudv1.Next)
				nextObject.Href = core.StringPtr("ibm.com?offset=135")
				responseObject.Next = nextObject

				value, err := responseObject.GetNextOffset()
				Expect(err).To(BeNil())
				Expect(value).To(Equal(core.Int64Ptr(int64(135))))
			})
			It(`Invoke GetNextOffset without a "Next" property in the response`, func() {
				responseObject := new(mqcloudv1.QueueManagerDetailsCollection)

				value, err := responseObject.GetNextOffset()
				Expect(err).To(BeNil())
				Expect(value).To(BeNil())
			})
			It(`Invoke GetNextOffset without any query params in the "Next" URL`, func() {
				responseObject := new(mqcloudv1.QueueManagerDetailsCollection)
				nextObject := new(mqcloudv1.Next)
				nextObject.Href = core.StringPtr("ibm.com")
				responseObject.Next = nextObject

				value, err := responseObject.GetNextOffset()
				Expect(err).To(BeNil())
				Expect(value).To(BeNil())
			})
			It(`Invoke GetNextOffset with a non-integer query param in the "Next" URL`, func() {
				responseObject := new(mqcloudv1.QueueManagerDetailsCollection)
				nextObject := new(mqcloudv1.Next)
				nextObject.Href = core.StringPtr("ibm.com?offset=tiger")
				responseObject.Next = nextObject

				value, err := responseObject.GetNextOffset()
				Expect(err).NotTo(BeNil())
				Expect(value).To(BeNil())
			})
		})
		Context(`Using mock server endpoint - paginated response`, func() {
			BeforeEach(func() {
				var requestNumber int = 0
				testServer = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
					defer GinkgoRecover()

					// Verify the contents of the request
					Expect(req.URL.EscapedPath()).To(Equal(listQueueManagersPath))
					Expect(req.Method).To(Equal("GET"))

					// Set mock response
					res.Header().Set("Content-type", "application/json")
					res.WriteHeader(200)
					requestNumber++
					if requestNumber == 1 {
						fmt.Fprintf(res, "%s", `{"next":{"href":"https://myhost.com/somePath?offset=1"},"total_count":2,"limit":1,"queue_managers":[{"id":"ID","name":"Name","display_name":"DisplayName","location":"reserved-eu-de-cluster-f884","size":"small","status_uri":"StatusURI","version":"9.3.2_2","web_console_url":"WebConsoleURL","rest_api_endpoint_url":"RestApiEndpointURL","administrator_api_endpoint_url":"AdministratorApiEndpointURL","connection_info_uri":"ConnectionInfoURI","date_created":"2020-01-13T15:39:35.000Z","upgrade_available":true,"available_upgrade_versions_uri":"AvailableUpgradeVersionsURI","href":"Href"}]}`)
					} else if requestNumber == 2 {
						fmt.Fprintf(res, "%s", `{"total_count":2,"limit":1,"queue_managers":[{"id":"ID","name":"Name","display_name":"DisplayName","location":"reserved-eu-de-cluster-f884","size":"small","status_uri":"StatusURI","version":"9.3.2_2","web_console_url":"WebConsoleURL","rest_api_endpoint_url":"RestApiEndpointURL","administrator_api_endpoint_url":"AdministratorApiEndpointURL","connection_info_uri":"ConnectionInfoURI","date_created":"2020-01-13T15:39:35.000Z","upgrade_available":true,"available_upgrade_versions_uri":"AvailableUpgradeVersionsURI","href":"Href"}]}`)
					} else {
						res.WriteHeader(400)
					}
				}))
			})
			It(`Use QueueManagersPager.GetNext successfully`, func() {
				mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1(&mqcloudv1.MqcloudV1Options{
					URL:            testServer.URL,
					Authenticator:  &core.NoAuthAuthenticator{},
					AcceptLanguage: core.StringPtr(acceptLanguage),
				})
				Expect(serviceErr).To(BeNil())
				Expect(mqcloudService).ToNot(BeNil())

				listQueueManagersOptionsModel := &mqcloudv1.ListQueueManagersOptions{
					ServiceInstanceGuid: core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8"),
					Limit:               core.Int64Ptr(int64(10)),
				}

				pager, err := mqcloudService.NewQueueManagersPager(listQueueManagersOptionsModel)
				Expect(err).To(BeNil())
				Expect(pager).ToNot(BeNil())

				var allResults []mqcloudv1.QueueManagerDetails
				for pager.HasNext() {
					nextPage, err := pager.GetNext()
					Expect(err).To(BeNil())
					Expect(nextPage).ToNot(BeNil())
					allResults = append(allResults, nextPage...)
				}
				Expect(len(allResults)).To(Equal(2))
			})
			It(`Use QueueManagersPager.GetAll successfully`, func() {
				mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1(&mqcloudv1.MqcloudV1Options{
					URL:            testServer.URL,
					Authenticator:  &core.NoAuthAuthenticator{},
					AcceptLanguage: core.StringPtr(acceptLanguage),
				})
				Expect(serviceErr).To(BeNil())
				Expect(mqcloudService).ToNot(BeNil())

				listQueueManagersOptionsModel := &mqcloudv1.ListQueueManagersOptions{
					ServiceInstanceGuid: core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8"),
					Limit:               core.Int64Ptr(int64(10)),
				}

				pager, err := mqcloudService.NewQueueManagersPager(listQueueManagersOptionsModel)
				Expect(err).To(BeNil())
				Expect(pager).ToNot(BeNil())

				allResults, err := pager.GetAll()
				Expect(err).To(BeNil())
				Expect(allResults).ToNot(BeNil())
				Expect(len(allResults)).To(Equal(2))
			})
		})
	})
	Describe(`GetQueueManager(getQueueManagerOptions *GetQueueManagerOptions) - Operation response error`, func() {
		acceptLanguage := "en-US,en;q=0.5"
		getQueueManagerPath := "/v1/a2b4d4bc-dadb-4637-bcec-9b7d1e723af8/queue_managers/b8e1aeda078009cf3db74e90d5d42328"
		Context(`Using mock server endpoint with invalid JSON response`, func() {
			BeforeEach(func() {
				testServer = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
					defer GinkgoRecover()

					// Verify the contents of the request
					Expect(req.URL.EscapedPath()).To(Equal(getQueueManagerPath))
					Expect(req.Method).To(Equal("GET"))
					Expect(req.Header["Accept-Language"]).ToNot(BeNil())
					Expect(req.Header["Accept-Language"][0]).To(Equal(fmt.Sprintf("%v", "en-US,en;q=0.5")))
					res.Header().Set("Content-type", "application/json")
					res.WriteHeader(200)
					fmt.Fprint(res, `} this is not valid json {`)
				}))
			})
			It(`Invoke GetQueueManager with error: Operation response processing error`, func() {
				mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1(&mqcloudv1.MqcloudV1Options{
					URL:            testServer.URL,
					Authenticator:  &core.NoAuthAuthenticator{},
					AcceptLanguage: core.StringPtr(acceptLanguage),
				})
				Expect(serviceErr).To(BeNil())
				Expect(mqcloudService).ToNot(BeNil())

				// Construct an instance of the GetQueueManagerOptions model
				getQueueManagerOptionsModel := new(mqcloudv1.GetQueueManagerOptions)
				getQueueManagerOptionsModel.ServiceInstanceGuid = core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				getQueueManagerOptionsModel.QueueManagerID = core.StringPtr("b8e1aeda078009cf3db74e90d5d42328")
				getQueueManagerOptionsModel.Headers = map[string]string{"x-custom-header": "x-custom-value"}
				// Expect response parsing to fail since we are receiving a text/plain response
				result, response, operationErr := mqcloudService.GetQueueManager(getQueueManagerOptionsModel)
				Expect(operationErr).ToNot(BeNil())
				Expect(response).ToNot(BeNil())
				Expect(result).To(BeNil())

				// Enable retries and test again
				mqcloudService.EnableRetries(0, 0)
				result, response, operationErr = mqcloudService.GetQueueManager(getQueueManagerOptionsModel)
				Expect(operationErr).ToNot(BeNil())
				Expect(response).ToNot(BeNil())
				Expect(result).To(BeNil())
			})
			AfterEach(func() {
				testServer.Close()
			})
		})
	})
	Describe(`GetQueueManager(getQueueManagerOptions *GetQueueManagerOptions)`, func() {
		acceptLanguage := "en-US,en;q=0.5"
		getQueueManagerPath := "/v1/a2b4d4bc-dadb-4637-bcec-9b7d1e723af8/queue_managers/b8e1aeda078009cf3db74e90d5d42328"
		Context(`Using mock server endpoint with timeout`, func() {
			BeforeEach(func() {
				testServer = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
					defer GinkgoRecover()

					// Verify the contents of the request
					Expect(req.URL.EscapedPath()).To(Equal(getQueueManagerPath))
					Expect(req.Method).To(Equal("GET"))

					Expect(req.Header["Accept-Language"]).ToNot(BeNil())
					Expect(req.Header["Accept-Language"][0]).To(Equal(fmt.Sprintf("%v", "en-US,en;q=0.5")))
					// Sleep a short time to support a timeout test
					time.Sleep(100 * time.Millisecond)

					// Set mock response
					res.Header().Set("Content-type", "application/json")
					res.WriteHeader(200)
					fmt.Fprintf(res, "%s", `{"id": "ID", "name": "Name", "display_name": "DisplayName", "location": "reserved-eu-de-cluster-f884", "size": "small", "status_uri": "StatusURI", "version": "9.3.2_2", "web_console_url": "WebConsoleURL", "rest_api_endpoint_url": "RestApiEndpointURL", "administrator_api_endpoint_url": "AdministratorApiEndpointURL", "connection_info_uri": "ConnectionInfoURI", "date_created": "2020-01-13T15:39:35.000Z", "upgrade_available": true, "available_upgrade_versions_uri": "AvailableUpgradeVersionsURI", "href": "Href"}`)
				}))
			})
			It(`Invoke GetQueueManager successfully with retries`, func() {
				mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1(&mqcloudv1.MqcloudV1Options{
					URL:            testServer.URL,
					Authenticator:  &core.NoAuthAuthenticator{},
					AcceptLanguage: core.StringPtr(acceptLanguage),
				})
				Expect(serviceErr).To(BeNil())
				Expect(mqcloudService).ToNot(BeNil())
				mqcloudService.EnableRetries(0, 0)

				// Construct an instance of the GetQueueManagerOptions model
				getQueueManagerOptionsModel := new(mqcloudv1.GetQueueManagerOptions)
				getQueueManagerOptionsModel.ServiceInstanceGuid = core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				getQueueManagerOptionsModel.QueueManagerID = core.StringPtr("b8e1aeda078009cf3db74e90d5d42328")
				getQueueManagerOptionsModel.Headers = map[string]string{"x-custom-header": "x-custom-value"}

				// Invoke operation with a Context to test a timeout error
				ctx, cancelFunc := context.WithTimeout(context.Background(), 80*time.Millisecond)
				defer cancelFunc()
				_, _, operationErr := mqcloudService.GetQueueManagerWithContext(ctx, getQueueManagerOptionsModel)
				Expect(operationErr).ToNot(BeNil())
				Expect(operationErr.Error()).To(ContainSubstring("deadline exceeded"))

				// Disable retries and test again
				mqcloudService.DisableRetries()
				result, response, operationErr := mqcloudService.GetQueueManager(getQueueManagerOptionsModel)
				Expect(operationErr).To(BeNil())
				Expect(response).ToNot(BeNil())
				Expect(result).ToNot(BeNil())

				// Re-test the timeout error with retries disabled
				ctx, cancelFunc2 := context.WithTimeout(context.Background(), 80*time.Millisecond)
				defer cancelFunc2()
				_, _, operationErr = mqcloudService.GetQueueManagerWithContext(ctx, getQueueManagerOptionsModel)
				Expect(operationErr).ToNot(BeNil())
				Expect(operationErr.Error()).To(ContainSubstring("deadline exceeded"))
			})
			AfterEach(func() {
				testServer.Close()
			})
		})
		Context(`Using mock server endpoint`, func() {
			BeforeEach(func() {
				testServer = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
					defer GinkgoRecover()

					// Verify the contents of the request
					Expect(req.URL.EscapedPath()).To(Equal(getQueueManagerPath))
					Expect(req.Method).To(Equal("GET"))

					Expect(req.Header["Accept-Language"]).ToNot(BeNil())
					Expect(req.Header["Accept-Language"][0]).To(Equal(fmt.Sprintf("%v", "en-US,en;q=0.5")))
					// Set mock response
					res.Header().Set("Content-type", "application/json")
					res.WriteHeader(200)
					fmt.Fprintf(res, "%s", `{"id": "ID", "name": "Name", "display_name": "DisplayName", "location": "reserved-eu-de-cluster-f884", "size": "small", "status_uri": "StatusURI", "version": "9.3.2_2", "web_console_url": "WebConsoleURL", "rest_api_endpoint_url": "RestApiEndpointURL", "administrator_api_endpoint_url": "AdministratorApiEndpointURL", "connection_info_uri": "ConnectionInfoURI", "date_created": "2020-01-13T15:39:35.000Z", "upgrade_available": true, "available_upgrade_versions_uri": "AvailableUpgradeVersionsURI", "href": "Href"}`)
				}))
			})
			It(`Invoke GetQueueManager successfully`, func() {
				mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1(&mqcloudv1.MqcloudV1Options{
					URL:            testServer.URL,
					Authenticator:  &core.NoAuthAuthenticator{},
					AcceptLanguage: core.StringPtr(acceptLanguage),
				})
				Expect(serviceErr).To(BeNil())
				Expect(mqcloudService).ToNot(BeNil())

				// Invoke operation with nil options model (negative test)
				result, response, operationErr := mqcloudService.GetQueueManager(nil)
				Expect(operationErr).NotTo(BeNil())
				Expect(response).To(BeNil())
				Expect(result).To(BeNil())

				// Construct an instance of the GetQueueManagerOptions model
				getQueueManagerOptionsModel := new(mqcloudv1.GetQueueManagerOptions)
				getQueueManagerOptionsModel.ServiceInstanceGuid = core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				getQueueManagerOptionsModel.QueueManagerID = core.StringPtr("b8e1aeda078009cf3db74e90d5d42328")
				getQueueManagerOptionsModel.Headers = map[string]string{"x-custom-header": "x-custom-value"}

				// Invoke operation with valid options model (positive test)
				result, response, operationErr = mqcloudService.GetQueueManager(getQueueManagerOptionsModel)
				Expect(operationErr).To(BeNil())
				Expect(response).ToNot(BeNil())
				Expect(result).ToNot(BeNil())

			})
			It(`Invoke GetQueueManager with error: Operation validation and request error`, func() {
				mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1(&mqcloudv1.MqcloudV1Options{
					URL:            testServer.URL,
					Authenticator:  &core.NoAuthAuthenticator{},
					AcceptLanguage: core.StringPtr(acceptLanguage),
				})
				Expect(serviceErr).To(BeNil())
				Expect(mqcloudService).ToNot(BeNil())

				// Construct an instance of the GetQueueManagerOptions model
				getQueueManagerOptionsModel := new(mqcloudv1.GetQueueManagerOptions)
				getQueueManagerOptionsModel.ServiceInstanceGuid = core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				getQueueManagerOptionsModel.QueueManagerID = core.StringPtr("b8e1aeda078009cf3db74e90d5d42328")
				getQueueManagerOptionsModel.Headers = map[string]string{"x-custom-header": "x-custom-value"}
				// Invoke operation with empty URL (negative test)
				err := mqcloudService.SetServiceURL("")
				Expect(err).To(BeNil())
				result, response, operationErr := mqcloudService.GetQueueManager(getQueueManagerOptionsModel)
				Expect(operationErr).ToNot(BeNil())
				Expect(operationErr.Error()).To(ContainSubstring(core.ERRORMSG_SERVICE_URL_MISSING))
				Expect(response).To(BeNil())
				Expect(result).To(BeNil())
				// Construct a second instance of the GetQueueManagerOptions model with no property values
				getQueueManagerOptionsModelNew := new(mqcloudv1.GetQueueManagerOptions)
				// Invoke operation with invalid model (negative test)
				result, response, operationErr = mqcloudService.GetQueueManager(getQueueManagerOptionsModelNew)
				Expect(operationErr).ToNot(BeNil())
				Expect(response).To(BeNil())
				Expect(result).To(BeNil())
			})
			AfterEach(func() {
				testServer.Close()
			})
		})
		Context(`Using mock server endpoint with missing response body`, func() {
			BeforeEach(func() {
				testServer = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
					defer GinkgoRecover()

					// Set success status code with no respoonse body
					res.WriteHeader(200)
				}))
			})
			It(`Invoke GetQueueManager successfully`, func() {
				mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1(&mqcloudv1.MqcloudV1Options{
					URL:            testServer.URL,
					Authenticator:  &core.NoAuthAuthenticator{},
					AcceptLanguage: core.StringPtr(acceptLanguage),
				})
				Expect(serviceErr).To(BeNil())
				Expect(mqcloudService).ToNot(BeNil())

				// Construct an instance of the GetQueueManagerOptions model
				getQueueManagerOptionsModel := new(mqcloudv1.GetQueueManagerOptions)
				getQueueManagerOptionsModel.ServiceInstanceGuid = core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				getQueueManagerOptionsModel.QueueManagerID = core.StringPtr("b8e1aeda078009cf3db74e90d5d42328")
				getQueueManagerOptionsModel.Headers = map[string]string{"x-custom-header": "x-custom-value"}

				// Invoke operation
				result, response, operationErr := mqcloudService.GetQueueManager(getQueueManagerOptionsModel)
				Expect(operationErr).To(BeNil())
				Expect(response).ToNot(BeNil())

				// Verify a nil result
				Expect(result).To(BeNil())
			})
			AfterEach(func() {
				testServer.Close()
			})
		})
	})
	Describe(`DeleteQueueManager(deleteQueueManagerOptions *DeleteQueueManagerOptions) - Operation response error`, func() {
		acceptLanguage := "en-US,en;q=0.5"
		deleteQueueManagerPath := "/v1/a2b4d4bc-dadb-4637-bcec-9b7d1e723af8/queue_managers/b8e1aeda078009cf3db74e90d5d42328"
		Context(`Using mock server endpoint with invalid JSON response`, func() {
			BeforeEach(func() {
				testServer = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
					defer GinkgoRecover()

					// Verify the contents of the request
					Expect(req.URL.EscapedPath()).To(Equal(deleteQueueManagerPath))
					Expect(req.Method).To(Equal("DELETE"))
					Expect(req.Header["Accept-Language"]).ToNot(BeNil())
					Expect(req.Header["Accept-Language"][0]).To(Equal(fmt.Sprintf("%v", "en-US,en;q=0.5")))
					res.Header().Set("Content-type", "application/json")
					res.WriteHeader(202)
					fmt.Fprint(res, `} this is not valid json {`)
				}))
			})
			It(`Invoke DeleteQueueManager with error: Operation response processing error`, func() {
				mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1(&mqcloudv1.MqcloudV1Options{
					URL:            testServer.URL,
					Authenticator:  &core.NoAuthAuthenticator{},
					AcceptLanguage: core.StringPtr(acceptLanguage),
				})
				Expect(serviceErr).To(BeNil())
				Expect(mqcloudService).ToNot(BeNil())

				// Construct an instance of the DeleteQueueManagerOptions model
				deleteQueueManagerOptionsModel := new(mqcloudv1.DeleteQueueManagerOptions)
				deleteQueueManagerOptionsModel.ServiceInstanceGuid = core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				deleteQueueManagerOptionsModel.QueueManagerID = core.StringPtr("b8e1aeda078009cf3db74e90d5d42328")
				deleteQueueManagerOptionsModel.Headers = map[string]string{"x-custom-header": "x-custom-value"}
				// Expect response parsing to fail since we are receiving a text/plain response
				result, response, operationErr := mqcloudService.DeleteQueueManager(deleteQueueManagerOptionsModel)
				Expect(operationErr).ToNot(BeNil())
				Expect(response).ToNot(BeNil())
				Expect(result).To(BeNil())

				// Enable retries and test again
				mqcloudService.EnableRetries(0, 0)
				result, response, operationErr = mqcloudService.DeleteQueueManager(deleteQueueManagerOptionsModel)
				Expect(operationErr).ToNot(BeNil())
				Expect(response).ToNot(BeNil())
				Expect(result).To(BeNil())
			})
			AfterEach(func() {
				testServer.Close()
			})
		})
	})
	Describe(`DeleteQueueManager(deleteQueueManagerOptions *DeleteQueueManagerOptions)`, func() {
		acceptLanguage := "en-US,en;q=0.5"
		deleteQueueManagerPath := "/v1/a2b4d4bc-dadb-4637-bcec-9b7d1e723af8/queue_managers/b8e1aeda078009cf3db74e90d5d42328"
		Context(`Using mock server endpoint with timeout`, func() {
			BeforeEach(func() {
				testServer = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
					defer GinkgoRecover()

					// Verify the contents of the request
					Expect(req.URL.EscapedPath()).To(Equal(deleteQueueManagerPath))
					Expect(req.Method).To(Equal("DELETE"))

					Expect(req.Header["Accept-Language"]).ToNot(BeNil())
					Expect(req.Header["Accept-Language"][0]).To(Equal(fmt.Sprintf("%v", "en-US,en;q=0.5")))
					// Sleep a short time to support a timeout test
					time.Sleep(100 * time.Millisecond)

					// Set mock response
					res.Header().Set("Content-type", "application/json")
					res.WriteHeader(202)
					fmt.Fprintf(res, "%s", `{"queue_manager_uri": "QueueManagerURI", "queue_manager_status_uri": "QueueManagerStatusURI", "queue_manager_id": "QueueManagerID"}`)
				}))
			})
			It(`Invoke DeleteQueueManager successfully with retries`, func() {
				mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1(&mqcloudv1.MqcloudV1Options{
					URL:            testServer.URL,
					Authenticator:  &core.NoAuthAuthenticator{},
					AcceptLanguage: core.StringPtr(acceptLanguage),
				})
				Expect(serviceErr).To(BeNil())
				Expect(mqcloudService).ToNot(BeNil())
				mqcloudService.EnableRetries(0, 0)

				// Construct an instance of the DeleteQueueManagerOptions model
				deleteQueueManagerOptionsModel := new(mqcloudv1.DeleteQueueManagerOptions)
				deleteQueueManagerOptionsModel.ServiceInstanceGuid = core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				deleteQueueManagerOptionsModel.QueueManagerID = core.StringPtr("b8e1aeda078009cf3db74e90d5d42328")
				deleteQueueManagerOptionsModel.Headers = map[string]string{"x-custom-header": "x-custom-value"}

				// Invoke operation with a Context to test a timeout error
				ctx, cancelFunc := context.WithTimeout(context.Background(), 80*time.Millisecond)
				defer cancelFunc()
				_, _, operationErr := mqcloudService.DeleteQueueManagerWithContext(ctx, deleteQueueManagerOptionsModel)
				Expect(operationErr).ToNot(BeNil())
				Expect(operationErr.Error()).To(ContainSubstring("deadline exceeded"))

				// Disable retries and test again
				mqcloudService.DisableRetries()
				result, response, operationErr := mqcloudService.DeleteQueueManager(deleteQueueManagerOptionsModel)
				Expect(operationErr).To(BeNil())
				Expect(response).ToNot(BeNil())
				Expect(result).ToNot(BeNil())

				// Re-test the timeout error with retries disabled
				ctx, cancelFunc2 := context.WithTimeout(context.Background(), 80*time.Millisecond)
				defer cancelFunc2()
				_, _, operationErr = mqcloudService.DeleteQueueManagerWithContext(ctx, deleteQueueManagerOptionsModel)
				Expect(operationErr).ToNot(BeNil())
				Expect(operationErr.Error()).To(ContainSubstring("deadline exceeded"))
			})
			AfterEach(func() {
				testServer.Close()
			})
		})
		Context(`Using mock server endpoint`, func() {
			BeforeEach(func() {
				testServer = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
					defer GinkgoRecover()

					// Verify the contents of the request
					Expect(req.URL.EscapedPath()).To(Equal(deleteQueueManagerPath))
					Expect(req.Method).To(Equal("DELETE"))

					Expect(req.Header["Accept-Language"]).ToNot(BeNil())
					Expect(req.Header["Accept-Language"][0]).To(Equal(fmt.Sprintf("%v", "en-US,en;q=0.5")))
					// Set mock response
					res.Header().Set("Content-type", "application/json")
					res.WriteHeader(202)
					fmt.Fprintf(res, "%s", `{"queue_manager_uri": "QueueManagerURI", "queue_manager_status_uri": "QueueManagerStatusURI", "queue_manager_id": "QueueManagerID"}`)
				}))
			})
			It(`Invoke DeleteQueueManager successfully`, func() {
				mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1(&mqcloudv1.MqcloudV1Options{
					URL:            testServer.URL,
					Authenticator:  &core.NoAuthAuthenticator{},
					AcceptLanguage: core.StringPtr(acceptLanguage),
				})
				Expect(serviceErr).To(BeNil())
				Expect(mqcloudService).ToNot(BeNil())

				// Invoke operation with nil options model (negative test)
				result, response, operationErr := mqcloudService.DeleteQueueManager(nil)
				Expect(operationErr).NotTo(BeNil())
				Expect(response).To(BeNil())
				Expect(result).To(BeNil())

				// Construct an instance of the DeleteQueueManagerOptions model
				deleteQueueManagerOptionsModel := new(mqcloudv1.DeleteQueueManagerOptions)
				deleteQueueManagerOptionsModel.ServiceInstanceGuid = core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				deleteQueueManagerOptionsModel.QueueManagerID = core.StringPtr("b8e1aeda078009cf3db74e90d5d42328")
				deleteQueueManagerOptionsModel.Headers = map[string]string{"x-custom-header": "x-custom-value"}

				// Invoke operation with valid options model (positive test)
				result, response, operationErr = mqcloudService.DeleteQueueManager(deleteQueueManagerOptionsModel)
				Expect(operationErr).To(BeNil())
				Expect(response).ToNot(BeNil())
				Expect(result).ToNot(BeNil())

			})
			It(`Invoke DeleteQueueManager with error: Operation validation and request error`, func() {
				mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1(&mqcloudv1.MqcloudV1Options{
					URL:            testServer.URL,
					Authenticator:  &core.NoAuthAuthenticator{},
					AcceptLanguage: core.StringPtr(acceptLanguage),
				})
				Expect(serviceErr).To(BeNil())
				Expect(mqcloudService).ToNot(BeNil())

				// Construct an instance of the DeleteQueueManagerOptions model
				deleteQueueManagerOptionsModel := new(mqcloudv1.DeleteQueueManagerOptions)
				deleteQueueManagerOptionsModel.ServiceInstanceGuid = core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				deleteQueueManagerOptionsModel.QueueManagerID = core.StringPtr("b8e1aeda078009cf3db74e90d5d42328")
				deleteQueueManagerOptionsModel.Headers = map[string]string{"x-custom-header": "x-custom-value"}
				// Invoke operation with empty URL (negative test)
				err := mqcloudService.SetServiceURL("")
				Expect(err).To(BeNil())
				result, response, operationErr := mqcloudService.DeleteQueueManager(deleteQueueManagerOptionsModel)
				Expect(operationErr).ToNot(BeNil())
				Expect(operationErr.Error()).To(ContainSubstring(core.ERRORMSG_SERVICE_URL_MISSING))
				Expect(response).To(BeNil())
				Expect(result).To(BeNil())
				// Construct a second instance of the DeleteQueueManagerOptions model with no property values
				deleteQueueManagerOptionsModelNew := new(mqcloudv1.DeleteQueueManagerOptions)
				// Invoke operation with invalid model (negative test)
				result, response, operationErr = mqcloudService.DeleteQueueManager(deleteQueueManagerOptionsModelNew)
				Expect(operationErr).ToNot(BeNil())
				Expect(response).To(BeNil())
				Expect(result).To(BeNil())
			})
			AfterEach(func() {
				testServer.Close()
			})
		})
		Context(`Using mock server endpoint with missing response body`, func() {
			BeforeEach(func() {
				testServer = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
					defer GinkgoRecover()

					// Set success status code with no respoonse body
					res.WriteHeader(202)
				}))
			})
			It(`Invoke DeleteQueueManager successfully`, func() {
				mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1(&mqcloudv1.MqcloudV1Options{
					URL:            testServer.URL,
					Authenticator:  &core.NoAuthAuthenticator{},
					AcceptLanguage: core.StringPtr(acceptLanguage),
				})
				Expect(serviceErr).To(BeNil())
				Expect(mqcloudService).ToNot(BeNil())

				// Construct an instance of the DeleteQueueManagerOptions model
				deleteQueueManagerOptionsModel := new(mqcloudv1.DeleteQueueManagerOptions)
				deleteQueueManagerOptionsModel.ServiceInstanceGuid = core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				deleteQueueManagerOptionsModel.QueueManagerID = core.StringPtr("b8e1aeda078009cf3db74e90d5d42328")
				deleteQueueManagerOptionsModel.Headers = map[string]string{"x-custom-header": "x-custom-value"}

				// Invoke operation
				result, response, operationErr := mqcloudService.DeleteQueueManager(deleteQueueManagerOptionsModel)
				Expect(operationErr).To(BeNil())
				Expect(response).ToNot(BeNil())

				// Verify a nil result
				Expect(result).To(BeNil())
			})
			AfterEach(func() {
				testServer.Close()
			})
		})
	})
	Describe(`SetQueueManagerVersion(setQueueManagerVersionOptions *SetQueueManagerVersionOptions) - Operation response error`, func() {
		acceptLanguage := "en-US,en;q=0.5"
		setQueueManagerVersionPath := "/v1/a2b4d4bc-dadb-4637-bcec-9b7d1e723af8/queue_managers/b8e1aeda078009cf3db74e90d5d42328/version"
		Context(`Using mock server endpoint with invalid JSON response`, func() {
			BeforeEach(func() {
				testServer = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
					defer GinkgoRecover()

					// Verify the contents of the request
					Expect(req.URL.EscapedPath()).To(Equal(setQueueManagerVersionPath))
					Expect(req.Method).To(Equal("PUT"))
					Expect(req.Header["Accept-Language"]).ToNot(BeNil())
					Expect(req.Header["Accept-Language"][0]).To(Equal(fmt.Sprintf("%v", "en-US,en;q=0.5")))
					res.Header().Set("Content-type", "application/json")
					res.WriteHeader(202)
					fmt.Fprint(res, `} this is not valid json {`)
				}))
			})
			It(`Invoke SetQueueManagerVersion with error: Operation response processing error`, func() {
				mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1(&mqcloudv1.MqcloudV1Options{
					URL:            testServer.URL,
					Authenticator:  &core.NoAuthAuthenticator{},
					AcceptLanguage: core.StringPtr(acceptLanguage),
				})
				Expect(serviceErr).To(BeNil())
				Expect(mqcloudService).ToNot(BeNil())

				// Construct an instance of the SetQueueManagerVersionOptions model
				setQueueManagerVersionOptionsModel := new(mqcloudv1.SetQueueManagerVersionOptions)
				setQueueManagerVersionOptionsModel.ServiceInstanceGuid = core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				setQueueManagerVersionOptionsModel.QueueManagerID = core.StringPtr("b8e1aeda078009cf3db74e90d5d42328")
				setQueueManagerVersionOptionsModel.Version = core.StringPtr("9.3.2_2")
				setQueueManagerVersionOptionsModel.Headers = map[string]string{"x-custom-header": "x-custom-value"}
				// Expect response parsing to fail since we are receiving a text/plain response
				result, response, operationErr := mqcloudService.SetQueueManagerVersion(setQueueManagerVersionOptionsModel)
				Expect(operationErr).ToNot(BeNil())
				Expect(response).ToNot(BeNil())
				Expect(result).To(BeNil())

				// Enable retries and test again
				mqcloudService.EnableRetries(0, 0)
				result, response, operationErr = mqcloudService.SetQueueManagerVersion(setQueueManagerVersionOptionsModel)
				Expect(operationErr).ToNot(BeNil())
				Expect(response).ToNot(BeNil())
				Expect(result).To(BeNil())
			})
			AfterEach(func() {
				testServer.Close()
			})
		})
	})
	Describe(`SetQueueManagerVersion(setQueueManagerVersionOptions *SetQueueManagerVersionOptions)`, func() {
		acceptLanguage := "en-US,en;q=0.5"
		setQueueManagerVersionPath := "/v1/a2b4d4bc-dadb-4637-bcec-9b7d1e723af8/queue_managers/b8e1aeda078009cf3db74e90d5d42328/version"
		Context(`Using mock server endpoint with timeout`, func() {
			BeforeEach(func() {
				testServer = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
					defer GinkgoRecover()

					// Verify the contents of the request
					Expect(req.URL.EscapedPath()).To(Equal(setQueueManagerVersionPath))
					Expect(req.Method).To(Equal("PUT"))

					// For gzip-disabled operation, verify Content-Encoding is not set.
					Expect(req.Header.Get("Content-Encoding")).To(BeEmpty())

					// If there is a body, then make sure we can read it
					bodyBuf := new(bytes.Buffer)
					if req.Header.Get("Content-Encoding") == "gzip" {
						body, err := core.NewGzipDecompressionReader(req.Body)
						Expect(err).To(BeNil())
						_, err = bodyBuf.ReadFrom(body)
						Expect(err).To(BeNil())
					} else {
						_, err := bodyBuf.ReadFrom(req.Body)
						Expect(err).To(BeNil())
					}
					fmt.Fprintf(GinkgoWriter, "  Request body: %s", bodyBuf.String())

					Expect(req.Header["Accept-Language"]).ToNot(BeNil())
					Expect(req.Header["Accept-Language"][0]).To(Equal(fmt.Sprintf("%v", "en-US,en;q=0.5")))
					// Sleep a short time to support a timeout test
					time.Sleep(100 * time.Millisecond)

					// Set mock response
					res.Header().Set("Content-type", "application/json")
					res.WriteHeader(202)
					fmt.Fprintf(res, "%s", `{"queue_manager_uri": "QueueManagerURI", "queue_manager_status_uri": "QueueManagerStatusURI", "queue_manager_id": "QueueManagerID"}`)
				}))
			})
			It(`Invoke SetQueueManagerVersion successfully with retries`, func() {
				mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1(&mqcloudv1.MqcloudV1Options{
					URL:            testServer.URL,
					Authenticator:  &core.NoAuthAuthenticator{},
					AcceptLanguage: core.StringPtr(acceptLanguage),
				})
				Expect(serviceErr).To(BeNil())
				Expect(mqcloudService).ToNot(BeNil())
				mqcloudService.EnableRetries(0, 0)

				// Construct an instance of the SetQueueManagerVersionOptions model
				setQueueManagerVersionOptionsModel := new(mqcloudv1.SetQueueManagerVersionOptions)
				setQueueManagerVersionOptionsModel.ServiceInstanceGuid = core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				setQueueManagerVersionOptionsModel.QueueManagerID = core.StringPtr("b8e1aeda078009cf3db74e90d5d42328")
				setQueueManagerVersionOptionsModel.Version = core.StringPtr("9.3.2_2")
				setQueueManagerVersionOptionsModel.Headers = map[string]string{"x-custom-header": "x-custom-value"}

				// Invoke operation with a Context to test a timeout error
				ctx, cancelFunc := context.WithTimeout(context.Background(), 80*time.Millisecond)
				defer cancelFunc()
				_, _, operationErr := mqcloudService.SetQueueManagerVersionWithContext(ctx, setQueueManagerVersionOptionsModel)
				Expect(operationErr).ToNot(BeNil())
				Expect(operationErr.Error()).To(ContainSubstring("deadline exceeded"))

				// Disable retries and test again
				mqcloudService.DisableRetries()
				result, response, operationErr := mqcloudService.SetQueueManagerVersion(setQueueManagerVersionOptionsModel)
				Expect(operationErr).To(BeNil())
				Expect(response).ToNot(BeNil())
				Expect(result).ToNot(BeNil())

				// Re-test the timeout error with retries disabled
				ctx, cancelFunc2 := context.WithTimeout(context.Background(), 80*time.Millisecond)
				defer cancelFunc2()
				_, _, operationErr = mqcloudService.SetQueueManagerVersionWithContext(ctx, setQueueManagerVersionOptionsModel)
				Expect(operationErr).ToNot(BeNil())
				Expect(operationErr.Error()).To(ContainSubstring("deadline exceeded"))
			})
			AfterEach(func() {
				testServer.Close()
			})
		})
		Context(`Using mock server endpoint`, func() {
			BeforeEach(func() {
				testServer = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
					defer GinkgoRecover()

					// Verify the contents of the request
					Expect(req.URL.EscapedPath()).To(Equal(setQueueManagerVersionPath))
					Expect(req.Method).To(Equal("PUT"))

					// For gzip-disabled operation, verify Content-Encoding is not set.
					Expect(req.Header.Get("Content-Encoding")).To(BeEmpty())

					// If there is a body, then make sure we can read it
					bodyBuf := new(bytes.Buffer)
					if req.Header.Get("Content-Encoding") == "gzip" {
						body, err := core.NewGzipDecompressionReader(req.Body)
						Expect(err).To(BeNil())
						_, err = bodyBuf.ReadFrom(body)
						Expect(err).To(BeNil())
					} else {
						_, err := bodyBuf.ReadFrom(req.Body)
						Expect(err).To(BeNil())
					}
					fmt.Fprintf(GinkgoWriter, "  Request body: %s", bodyBuf.String())

					Expect(req.Header["Accept-Language"]).ToNot(BeNil())
					Expect(req.Header["Accept-Language"][0]).To(Equal(fmt.Sprintf("%v", "en-US,en;q=0.5")))
					// Set mock response
					res.Header().Set("Content-type", "application/json")
					res.WriteHeader(202)
					fmt.Fprintf(res, "%s", `{"queue_manager_uri": "QueueManagerURI", "queue_manager_status_uri": "QueueManagerStatusURI", "queue_manager_id": "QueueManagerID"}`)
				}))
			})
			It(`Invoke SetQueueManagerVersion successfully`, func() {
				mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1(&mqcloudv1.MqcloudV1Options{
					URL:            testServer.URL,
					Authenticator:  &core.NoAuthAuthenticator{},
					AcceptLanguage: core.StringPtr(acceptLanguage),
				})
				Expect(serviceErr).To(BeNil())
				Expect(mqcloudService).ToNot(BeNil())

				// Invoke operation with nil options model (negative test)
				result, response, operationErr := mqcloudService.SetQueueManagerVersion(nil)
				Expect(operationErr).NotTo(BeNil())
				Expect(response).To(BeNil())
				Expect(result).To(BeNil())

				// Construct an instance of the SetQueueManagerVersionOptions model
				setQueueManagerVersionOptionsModel := new(mqcloudv1.SetQueueManagerVersionOptions)
				setQueueManagerVersionOptionsModel.ServiceInstanceGuid = core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				setQueueManagerVersionOptionsModel.QueueManagerID = core.StringPtr("b8e1aeda078009cf3db74e90d5d42328")
				setQueueManagerVersionOptionsModel.Version = core.StringPtr("9.3.2_2")
				setQueueManagerVersionOptionsModel.Headers = map[string]string{"x-custom-header": "x-custom-value"}

				// Invoke operation with valid options model (positive test)
				result, response, operationErr = mqcloudService.SetQueueManagerVersion(setQueueManagerVersionOptionsModel)
				Expect(operationErr).To(BeNil())
				Expect(response).ToNot(BeNil())
				Expect(result).ToNot(BeNil())

			})
			It(`Invoke SetQueueManagerVersion with error: Operation validation and request error`, func() {
				mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1(&mqcloudv1.MqcloudV1Options{
					URL:            testServer.URL,
					Authenticator:  &core.NoAuthAuthenticator{},
					AcceptLanguage: core.StringPtr(acceptLanguage),
				})
				Expect(serviceErr).To(BeNil())
				Expect(mqcloudService).ToNot(BeNil())

				// Construct an instance of the SetQueueManagerVersionOptions model
				setQueueManagerVersionOptionsModel := new(mqcloudv1.SetQueueManagerVersionOptions)
				setQueueManagerVersionOptionsModel.ServiceInstanceGuid = core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				setQueueManagerVersionOptionsModel.QueueManagerID = core.StringPtr("b8e1aeda078009cf3db74e90d5d42328")
				setQueueManagerVersionOptionsModel.Version = core.StringPtr("9.3.2_2")
				setQueueManagerVersionOptionsModel.Headers = map[string]string{"x-custom-header": "x-custom-value"}
				// Invoke operation with empty URL (negative test)
				err := mqcloudService.SetServiceURL("")
				Expect(err).To(BeNil())
				result, response, operationErr := mqcloudService.SetQueueManagerVersion(setQueueManagerVersionOptionsModel)
				Expect(operationErr).ToNot(BeNil())
				Expect(operationErr.Error()).To(ContainSubstring(core.ERRORMSG_SERVICE_URL_MISSING))
				Expect(response).To(BeNil())
				Expect(result).To(BeNil())
				// Construct a second instance of the SetQueueManagerVersionOptions model with no property values
				setQueueManagerVersionOptionsModelNew := new(mqcloudv1.SetQueueManagerVersionOptions)
				// Invoke operation with invalid model (negative test)
				result, response, operationErr = mqcloudService.SetQueueManagerVersion(setQueueManagerVersionOptionsModelNew)
				Expect(operationErr).ToNot(BeNil())
				Expect(response).To(BeNil())
				Expect(result).To(BeNil())
			})
			AfterEach(func() {
				testServer.Close()
			})
		})
		Context(`Using mock server endpoint with missing response body`, func() {
			BeforeEach(func() {
				testServer = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
					defer GinkgoRecover()

					// Set success status code with no respoonse body
					res.WriteHeader(202)
				}))
			})
			It(`Invoke SetQueueManagerVersion successfully`, func() {
				mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1(&mqcloudv1.MqcloudV1Options{
					URL:            testServer.URL,
					Authenticator:  &core.NoAuthAuthenticator{},
					AcceptLanguage: core.StringPtr(acceptLanguage),
				})
				Expect(serviceErr).To(BeNil())
				Expect(mqcloudService).ToNot(BeNil())

				// Construct an instance of the SetQueueManagerVersionOptions model
				setQueueManagerVersionOptionsModel := new(mqcloudv1.SetQueueManagerVersionOptions)
				setQueueManagerVersionOptionsModel.ServiceInstanceGuid = core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				setQueueManagerVersionOptionsModel.QueueManagerID = core.StringPtr("b8e1aeda078009cf3db74e90d5d42328")
				setQueueManagerVersionOptionsModel.Version = core.StringPtr("9.3.2_2")
				setQueueManagerVersionOptionsModel.Headers = map[string]string{"x-custom-header": "x-custom-value"}

				// Invoke operation
				result, response, operationErr := mqcloudService.SetQueueManagerVersion(setQueueManagerVersionOptionsModel)
				Expect(operationErr).To(BeNil())
				Expect(response).ToNot(BeNil())

				// Verify a nil result
				Expect(result).To(BeNil())
			})
			AfterEach(func() {
				testServer.Close()
			})
		})
	})
	Describe(`GetQueueManagerAvailableUpgradeVersions(getQueueManagerAvailableUpgradeVersionsOptions *GetQueueManagerAvailableUpgradeVersionsOptions) - Operation response error`, func() {
		acceptLanguage := "en-US,en;q=0.5"
		getQueueManagerAvailableUpgradeVersionsPath := "/v1/a2b4d4bc-dadb-4637-bcec-9b7d1e723af8/queue_managers/b8e1aeda078009cf3db74e90d5d42328/available_versions"
		Context(`Using mock server endpoint with invalid JSON response`, func() {
			BeforeEach(func() {
				testServer = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
					defer GinkgoRecover()

					// Verify the contents of the request
					Expect(req.URL.EscapedPath()).To(Equal(getQueueManagerAvailableUpgradeVersionsPath))
					Expect(req.Method).To(Equal("GET"))
					Expect(req.Header["Accept-Language"]).ToNot(BeNil())
					Expect(req.Header["Accept-Language"][0]).To(Equal(fmt.Sprintf("%v", "en-US,en;q=0.5")))
					res.Header().Set("Content-type", "application/json")
					res.WriteHeader(200)
					fmt.Fprint(res, `} this is not valid json {`)
				}))
			})
			It(`Invoke GetQueueManagerAvailableUpgradeVersions with error: Operation response processing error`, func() {
				mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1(&mqcloudv1.MqcloudV1Options{
					URL:            testServer.URL,
					Authenticator:  &core.NoAuthAuthenticator{},
					AcceptLanguage: core.StringPtr(acceptLanguage),
				})
				Expect(serviceErr).To(BeNil())
				Expect(mqcloudService).ToNot(BeNil())

				// Construct an instance of the GetQueueManagerAvailableUpgradeVersionsOptions model
				getQueueManagerAvailableUpgradeVersionsOptionsModel := new(mqcloudv1.GetQueueManagerAvailableUpgradeVersionsOptions)
				getQueueManagerAvailableUpgradeVersionsOptionsModel.ServiceInstanceGuid = core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				getQueueManagerAvailableUpgradeVersionsOptionsModel.QueueManagerID = core.StringPtr("b8e1aeda078009cf3db74e90d5d42328")
				getQueueManagerAvailableUpgradeVersionsOptionsModel.Headers = map[string]string{"x-custom-header": "x-custom-value"}
				// Expect response parsing to fail since we are receiving a text/plain response
				result, response, operationErr := mqcloudService.GetQueueManagerAvailableUpgradeVersions(getQueueManagerAvailableUpgradeVersionsOptionsModel)
				Expect(operationErr).ToNot(BeNil())
				Expect(response).ToNot(BeNil())
				Expect(result).To(BeNil())

				// Enable retries and test again
				mqcloudService.EnableRetries(0, 0)
				result, response, operationErr = mqcloudService.GetQueueManagerAvailableUpgradeVersions(getQueueManagerAvailableUpgradeVersionsOptionsModel)
				Expect(operationErr).ToNot(BeNil())
				Expect(response).ToNot(BeNil())
				Expect(result).To(BeNil())
			})
			AfterEach(func() {
				testServer.Close()
			})
		})
	})
	Describe(`GetQueueManagerAvailableUpgradeVersions(getQueueManagerAvailableUpgradeVersionsOptions *GetQueueManagerAvailableUpgradeVersionsOptions)`, func() {
		acceptLanguage := "en-US,en;q=0.5"
		getQueueManagerAvailableUpgradeVersionsPath := "/v1/a2b4d4bc-dadb-4637-bcec-9b7d1e723af8/queue_managers/b8e1aeda078009cf3db74e90d5d42328/available_versions"
		Context(`Using mock server endpoint with timeout`, func() {
			BeforeEach(func() {
				testServer = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
					defer GinkgoRecover()

					// Verify the contents of the request
					Expect(req.URL.EscapedPath()).To(Equal(getQueueManagerAvailableUpgradeVersionsPath))
					Expect(req.Method).To(Equal("GET"))

					Expect(req.Header["Accept-Language"]).ToNot(BeNil())
					Expect(req.Header["Accept-Language"][0]).To(Equal(fmt.Sprintf("%v", "en-US,en;q=0.5")))
					// Sleep a short time to support a timeout test
					time.Sleep(100 * time.Millisecond)

					// Set mock response
					res.Header().Set("Content-type", "application/json")
					res.WriteHeader(200)
					fmt.Fprintf(res, "%s", `{"total_count": 10, "versions": [{"version": "9.3.2_2", "target_date": "2023-01-13T15:39:35.000Z"}]}`)
				}))
			})
			It(`Invoke GetQueueManagerAvailableUpgradeVersions successfully with retries`, func() {
				mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1(&mqcloudv1.MqcloudV1Options{
					URL:            testServer.URL,
					Authenticator:  &core.NoAuthAuthenticator{},
					AcceptLanguage: core.StringPtr(acceptLanguage),
				})
				Expect(serviceErr).To(BeNil())
				Expect(mqcloudService).ToNot(BeNil())
				mqcloudService.EnableRetries(0, 0)

				// Construct an instance of the GetQueueManagerAvailableUpgradeVersionsOptions model
				getQueueManagerAvailableUpgradeVersionsOptionsModel := new(mqcloudv1.GetQueueManagerAvailableUpgradeVersionsOptions)
				getQueueManagerAvailableUpgradeVersionsOptionsModel.ServiceInstanceGuid = core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				getQueueManagerAvailableUpgradeVersionsOptionsModel.QueueManagerID = core.StringPtr("b8e1aeda078009cf3db74e90d5d42328")
				getQueueManagerAvailableUpgradeVersionsOptionsModel.Headers = map[string]string{"x-custom-header": "x-custom-value"}

				// Invoke operation with a Context to test a timeout error
				ctx, cancelFunc := context.WithTimeout(context.Background(), 80*time.Millisecond)
				defer cancelFunc()
				_, _, operationErr := mqcloudService.GetQueueManagerAvailableUpgradeVersionsWithContext(ctx, getQueueManagerAvailableUpgradeVersionsOptionsModel)
				Expect(operationErr).ToNot(BeNil())
				Expect(operationErr.Error()).To(ContainSubstring("deadline exceeded"))

				// Disable retries and test again
				mqcloudService.DisableRetries()
				result, response, operationErr := mqcloudService.GetQueueManagerAvailableUpgradeVersions(getQueueManagerAvailableUpgradeVersionsOptionsModel)
				Expect(operationErr).To(BeNil())
				Expect(response).ToNot(BeNil())
				Expect(result).ToNot(BeNil())

				// Re-test the timeout error with retries disabled
				ctx, cancelFunc2 := context.WithTimeout(context.Background(), 80*time.Millisecond)
				defer cancelFunc2()
				_, _, operationErr = mqcloudService.GetQueueManagerAvailableUpgradeVersionsWithContext(ctx, getQueueManagerAvailableUpgradeVersionsOptionsModel)
				Expect(operationErr).ToNot(BeNil())
				Expect(operationErr.Error()).To(ContainSubstring("deadline exceeded"))
			})
			AfterEach(func() {
				testServer.Close()
			})
		})
		Context(`Using mock server endpoint`, func() {
			BeforeEach(func() {
				testServer = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
					defer GinkgoRecover()

					// Verify the contents of the request
					Expect(req.URL.EscapedPath()).To(Equal(getQueueManagerAvailableUpgradeVersionsPath))
					Expect(req.Method).To(Equal("GET"))

					Expect(req.Header["Accept-Language"]).ToNot(BeNil())
					Expect(req.Header["Accept-Language"][0]).To(Equal(fmt.Sprintf("%v", "en-US,en;q=0.5")))
					// Set mock response
					res.Header().Set("Content-type", "application/json")
					res.WriteHeader(200)
					fmt.Fprintf(res, "%s", `{"total_count": 10, "versions": [{"version": "9.3.2_2", "target_date": "2023-01-13T15:39:35.000Z"}]}`)
				}))
			})
			It(`Invoke GetQueueManagerAvailableUpgradeVersions successfully`, func() {
				mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1(&mqcloudv1.MqcloudV1Options{
					URL:            testServer.URL,
					Authenticator:  &core.NoAuthAuthenticator{},
					AcceptLanguage: core.StringPtr(acceptLanguage),
				})
				Expect(serviceErr).To(BeNil())
				Expect(mqcloudService).ToNot(BeNil())

				// Invoke operation with nil options model (negative test)
				result, response, operationErr := mqcloudService.GetQueueManagerAvailableUpgradeVersions(nil)
				Expect(operationErr).NotTo(BeNil())
				Expect(response).To(BeNil())
				Expect(result).To(BeNil())

				// Construct an instance of the GetQueueManagerAvailableUpgradeVersionsOptions model
				getQueueManagerAvailableUpgradeVersionsOptionsModel := new(mqcloudv1.GetQueueManagerAvailableUpgradeVersionsOptions)
				getQueueManagerAvailableUpgradeVersionsOptionsModel.ServiceInstanceGuid = core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				getQueueManagerAvailableUpgradeVersionsOptionsModel.QueueManagerID = core.StringPtr("b8e1aeda078009cf3db74e90d5d42328")
				getQueueManagerAvailableUpgradeVersionsOptionsModel.Headers = map[string]string{"x-custom-header": "x-custom-value"}

				// Invoke operation with valid options model (positive test)
				result, response, operationErr = mqcloudService.GetQueueManagerAvailableUpgradeVersions(getQueueManagerAvailableUpgradeVersionsOptionsModel)
				Expect(operationErr).To(BeNil())
				Expect(response).ToNot(BeNil())
				Expect(result).ToNot(BeNil())

			})
			It(`Invoke GetQueueManagerAvailableUpgradeVersions with error: Operation validation and request error`, func() {
				mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1(&mqcloudv1.MqcloudV1Options{
					URL:            testServer.URL,
					Authenticator:  &core.NoAuthAuthenticator{},
					AcceptLanguage: core.StringPtr(acceptLanguage),
				})
				Expect(serviceErr).To(BeNil())
				Expect(mqcloudService).ToNot(BeNil())

				// Construct an instance of the GetQueueManagerAvailableUpgradeVersionsOptions model
				getQueueManagerAvailableUpgradeVersionsOptionsModel := new(mqcloudv1.GetQueueManagerAvailableUpgradeVersionsOptions)
				getQueueManagerAvailableUpgradeVersionsOptionsModel.ServiceInstanceGuid = core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				getQueueManagerAvailableUpgradeVersionsOptionsModel.QueueManagerID = core.StringPtr("b8e1aeda078009cf3db74e90d5d42328")
				getQueueManagerAvailableUpgradeVersionsOptionsModel.Headers = map[string]string{"x-custom-header": "x-custom-value"}
				// Invoke operation with empty URL (negative test)
				err := mqcloudService.SetServiceURL("")
				Expect(err).To(BeNil())
				result, response, operationErr := mqcloudService.GetQueueManagerAvailableUpgradeVersions(getQueueManagerAvailableUpgradeVersionsOptionsModel)
				Expect(operationErr).ToNot(BeNil())
				Expect(operationErr.Error()).To(ContainSubstring(core.ERRORMSG_SERVICE_URL_MISSING))
				Expect(response).To(BeNil())
				Expect(result).To(BeNil())
				// Construct a second instance of the GetQueueManagerAvailableUpgradeVersionsOptions model with no property values
				getQueueManagerAvailableUpgradeVersionsOptionsModelNew := new(mqcloudv1.GetQueueManagerAvailableUpgradeVersionsOptions)
				// Invoke operation with invalid model (negative test)
				result, response, operationErr = mqcloudService.GetQueueManagerAvailableUpgradeVersions(getQueueManagerAvailableUpgradeVersionsOptionsModelNew)
				Expect(operationErr).ToNot(BeNil())
				Expect(response).To(BeNil())
				Expect(result).To(BeNil())
			})
			AfterEach(func() {
				testServer.Close()
			})
		})
		Context(`Using mock server endpoint with missing response body`, func() {
			BeforeEach(func() {
				testServer = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
					defer GinkgoRecover()

					// Set success status code with no respoonse body
					res.WriteHeader(200)
				}))
			})
			It(`Invoke GetQueueManagerAvailableUpgradeVersions successfully`, func() {
				mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1(&mqcloudv1.MqcloudV1Options{
					URL:            testServer.URL,
					Authenticator:  &core.NoAuthAuthenticator{},
					AcceptLanguage: core.StringPtr(acceptLanguage),
				})
				Expect(serviceErr).To(BeNil())
				Expect(mqcloudService).ToNot(BeNil())

				// Construct an instance of the GetQueueManagerAvailableUpgradeVersionsOptions model
				getQueueManagerAvailableUpgradeVersionsOptionsModel := new(mqcloudv1.GetQueueManagerAvailableUpgradeVersionsOptions)
				getQueueManagerAvailableUpgradeVersionsOptionsModel.ServiceInstanceGuid = core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				getQueueManagerAvailableUpgradeVersionsOptionsModel.QueueManagerID = core.StringPtr("b8e1aeda078009cf3db74e90d5d42328")
				getQueueManagerAvailableUpgradeVersionsOptionsModel.Headers = map[string]string{"x-custom-header": "x-custom-value"}

				// Invoke operation
				result, response, operationErr := mqcloudService.GetQueueManagerAvailableUpgradeVersions(getQueueManagerAvailableUpgradeVersionsOptionsModel)
				Expect(operationErr).To(BeNil())
				Expect(response).ToNot(BeNil())

				// Verify a nil result
				Expect(result).To(BeNil())
			})
			AfterEach(func() {
				testServer.Close()
			})
		})
	})
	Describe(`GetQueueManagerConnectionInfo(getQueueManagerConnectionInfoOptions *GetQueueManagerConnectionInfoOptions) - Operation response error`, func() {
		acceptLanguage := "en-US,en;q=0.5"
		getQueueManagerConnectionInfoPath := "/v1/a2b4d4bc-dadb-4637-bcec-9b7d1e723af8/queue_managers/b8e1aeda078009cf3db74e90d5d42328/connection_info"
		Context(`Using mock server endpoint with invalid JSON response`, func() {
			BeforeEach(func() {
				testServer = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
					defer GinkgoRecover()

					// Verify the contents of the request
					Expect(req.URL.EscapedPath()).To(Equal(getQueueManagerConnectionInfoPath))
					Expect(req.Method).To(Equal("GET"))
					Expect(req.Header["Accept-Language"]).ToNot(BeNil())
					Expect(req.Header["Accept-Language"][0]).To(Equal(fmt.Sprintf("%v", "en-US,en;q=0.5")))
					res.Header().Set("Content-type", "application/json")
					res.WriteHeader(200)
					fmt.Fprint(res, `} this is not valid json {`)
				}))
			})
			It(`Invoke GetQueueManagerConnectionInfo with error: Operation response processing error`, func() {
				mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1(&mqcloudv1.MqcloudV1Options{
					URL:            testServer.URL,
					Authenticator:  &core.NoAuthAuthenticator{},
					AcceptLanguage: core.StringPtr(acceptLanguage),
				})
				Expect(serviceErr).To(BeNil())
				Expect(mqcloudService).ToNot(BeNil())

				// Construct an instance of the GetQueueManagerConnectionInfoOptions model
				getQueueManagerConnectionInfoOptionsModel := new(mqcloudv1.GetQueueManagerConnectionInfoOptions)
				getQueueManagerConnectionInfoOptionsModel.ServiceInstanceGuid = core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				getQueueManagerConnectionInfoOptionsModel.QueueManagerID = core.StringPtr("b8e1aeda078009cf3db74e90d5d42328")
				getQueueManagerConnectionInfoOptionsModel.Headers = map[string]string{"x-custom-header": "x-custom-value"}
				// Expect response parsing to fail since we are receiving a text/plain response
				result, response, operationErr := mqcloudService.GetQueueManagerConnectionInfo(getQueueManagerConnectionInfoOptionsModel)
				Expect(operationErr).ToNot(BeNil())
				Expect(response).ToNot(BeNil())
				Expect(result).To(BeNil())

				// Enable retries and test again
				mqcloudService.EnableRetries(0, 0)
				result, response, operationErr = mqcloudService.GetQueueManagerConnectionInfo(getQueueManagerConnectionInfoOptionsModel)
				Expect(operationErr).ToNot(BeNil())
				Expect(response).ToNot(BeNil())
				Expect(result).To(BeNil())
			})
			AfterEach(func() {
				testServer.Close()
			})
		})
	})
	Describe(`GetQueueManagerConnectionInfo(getQueueManagerConnectionInfoOptions *GetQueueManagerConnectionInfoOptions)`, func() {
		acceptLanguage := "en-US,en;q=0.5"
		getQueueManagerConnectionInfoPath := "/v1/a2b4d4bc-dadb-4637-bcec-9b7d1e723af8/queue_managers/b8e1aeda078009cf3db74e90d5d42328/connection_info"
		Context(`Using mock server endpoint with timeout`, func() {
			BeforeEach(func() {
				testServer = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
					defer GinkgoRecover()

					// Verify the contents of the request
					Expect(req.URL.EscapedPath()).To(Equal(getQueueManagerConnectionInfoPath))
					Expect(req.Method).To(Equal("GET"))

					Expect(req.Header["Accept-Language"]).ToNot(BeNil())
					Expect(req.Header["Accept-Language"][0]).To(Equal(fmt.Sprintf("%v", "en-US,en;q=0.5")))
					// Sleep a short time to support a timeout test
					time.Sleep(100 * time.Millisecond)

					// Set mock response
					res.Header().Set("Content-type", "application/json")
					res.WriteHeader(200)
					fmt.Fprintf(res, "%s", `{"channel": [{"name": "Name", "clientConnection": {"connection": [{"host": "Host", "port": 4}], "queueManager": "QueueManager"}, "transmissionSecurity": {"cipherSpecification": "CipherSpecification"}, "type": "Type"}]}`)
				}))
			})
			It(`Invoke GetQueueManagerConnectionInfo successfully with retries`, func() {
				mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1(&mqcloudv1.MqcloudV1Options{
					URL:            testServer.URL,
					Authenticator:  &core.NoAuthAuthenticator{},
					AcceptLanguage: core.StringPtr(acceptLanguage),
				})
				Expect(serviceErr).To(BeNil())
				Expect(mqcloudService).ToNot(BeNil())
				mqcloudService.EnableRetries(0, 0)

				// Construct an instance of the GetQueueManagerConnectionInfoOptions model
				getQueueManagerConnectionInfoOptionsModel := new(mqcloudv1.GetQueueManagerConnectionInfoOptions)
				getQueueManagerConnectionInfoOptionsModel.ServiceInstanceGuid = core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				getQueueManagerConnectionInfoOptionsModel.QueueManagerID = core.StringPtr("b8e1aeda078009cf3db74e90d5d42328")
				getQueueManagerConnectionInfoOptionsModel.Headers = map[string]string{"x-custom-header": "x-custom-value"}

				// Invoke operation with a Context to test a timeout error
				ctx, cancelFunc := context.WithTimeout(context.Background(), 80*time.Millisecond)
				defer cancelFunc()
				_, _, operationErr := mqcloudService.GetQueueManagerConnectionInfoWithContext(ctx, getQueueManagerConnectionInfoOptionsModel)
				Expect(operationErr).ToNot(BeNil())
				Expect(operationErr.Error()).To(ContainSubstring("deadline exceeded"))

				// Disable retries and test again
				mqcloudService.DisableRetries()
				result, response, operationErr := mqcloudService.GetQueueManagerConnectionInfo(getQueueManagerConnectionInfoOptionsModel)
				Expect(operationErr).To(BeNil())
				Expect(response).ToNot(BeNil())
				Expect(result).ToNot(BeNil())

				// Re-test the timeout error with retries disabled
				ctx, cancelFunc2 := context.WithTimeout(context.Background(), 80*time.Millisecond)
				defer cancelFunc2()
				_, _, operationErr = mqcloudService.GetQueueManagerConnectionInfoWithContext(ctx, getQueueManagerConnectionInfoOptionsModel)
				Expect(operationErr).ToNot(BeNil())
				Expect(operationErr.Error()).To(ContainSubstring("deadline exceeded"))
			})
			AfterEach(func() {
				testServer.Close()
			})
		})
		Context(`Using mock server endpoint`, func() {
			BeforeEach(func() {
				testServer = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
					defer GinkgoRecover()

					// Verify the contents of the request
					Expect(req.URL.EscapedPath()).To(Equal(getQueueManagerConnectionInfoPath))
					Expect(req.Method).To(Equal("GET"))

					Expect(req.Header["Accept-Language"]).ToNot(BeNil())
					Expect(req.Header["Accept-Language"][0]).To(Equal(fmt.Sprintf("%v", "en-US,en;q=0.5")))
					// Set mock response
					res.Header().Set("Content-type", "application/json")
					res.WriteHeader(200)
					fmt.Fprintf(res, "%s", `{"channel": [{"name": "Name", "clientConnection": {"connection": [{"host": "Host", "port": 4}], "queueManager": "QueueManager"}, "transmissionSecurity": {"cipherSpecification": "CipherSpecification"}, "type": "Type"}]}`)
				}))
			})
			It(`Invoke GetQueueManagerConnectionInfo successfully`, func() {
				mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1(&mqcloudv1.MqcloudV1Options{
					URL:            testServer.URL,
					Authenticator:  &core.NoAuthAuthenticator{},
					AcceptLanguage: core.StringPtr(acceptLanguage),
				})
				Expect(serviceErr).To(BeNil())
				Expect(mqcloudService).ToNot(BeNil())

				// Invoke operation with nil options model (negative test)
				result, response, operationErr := mqcloudService.GetQueueManagerConnectionInfo(nil)
				Expect(operationErr).NotTo(BeNil())
				Expect(response).To(BeNil())
				Expect(result).To(BeNil())

				// Construct an instance of the GetQueueManagerConnectionInfoOptions model
				getQueueManagerConnectionInfoOptionsModel := new(mqcloudv1.GetQueueManagerConnectionInfoOptions)
				getQueueManagerConnectionInfoOptionsModel.ServiceInstanceGuid = core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				getQueueManagerConnectionInfoOptionsModel.QueueManagerID = core.StringPtr("b8e1aeda078009cf3db74e90d5d42328")
				getQueueManagerConnectionInfoOptionsModel.Headers = map[string]string{"x-custom-header": "x-custom-value"}

				// Invoke operation with valid options model (positive test)
				result, response, operationErr = mqcloudService.GetQueueManagerConnectionInfo(getQueueManagerConnectionInfoOptionsModel)
				Expect(operationErr).To(BeNil())
				Expect(response).ToNot(BeNil())
				Expect(result).ToNot(BeNil())

			})
			It(`Invoke GetQueueManagerConnectionInfo with error: Operation validation and request error`, func() {
				mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1(&mqcloudv1.MqcloudV1Options{
					URL:            testServer.URL,
					Authenticator:  &core.NoAuthAuthenticator{},
					AcceptLanguage: core.StringPtr(acceptLanguage),
				})
				Expect(serviceErr).To(BeNil())
				Expect(mqcloudService).ToNot(BeNil())

				// Construct an instance of the GetQueueManagerConnectionInfoOptions model
				getQueueManagerConnectionInfoOptionsModel := new(mqcloudv1.GetQueueManagerConnectionInfoOptions)
				getQueueManagerConnectionInfoOptionsModel.ServiceInstanceGuid = core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				getQueueManagerConnectionInfoOptionsModel.QueueManagerID = core.StringPtr("b8e1aeda078009cf3db74e90d5d42328")
				getQueueManagerConnectionInfoOptionsModel.Headers = map[string]string{"x-custom-header": "x-custom-value"}
				// Invoke operation with empty URL (negative test)
				err := mqcloudService.SetServiceURL("")
				Expect(err).To(BeNil())
				result, response, operationErr := mqcloudService.GetQueueManagerConnectionInfo(getQueueManagerConnectionInfoOptionsModel)
				Expect(operationErr).ToNot(BeNil())
				Expect(operationErr.Error()).To(ContainSubstring(core.ERRORMSG_SERVICE_URL_MISSING))
				Expect(response).To(BeNil())
				Expect(result).To(BeNil())
				// Construct a second instance of the GetQueueManagerConnectionInfoOptions model with no property values
				getQueueManagerConnectionInfoOptionsModelNew := new(mqcloudv1.GetQueueManagerConnectionInfoOptions)
				// Invoke operation with invalid model (negative test)
				result, response, operationErr = mqcloudService.GetQueueManagerConnectionInfo(getQueueManagerConnectionInfoOptionsModelNew)
				Expect(operationErr).ToNot(BeNil())
				Expect(response).To(BeNil())
				Expect(result).To(BeNil())
			})
			AfterEach(func() {
				testServer.Close()
			})
		})
		Context(`Using mock server endpoint with missing response body`, func() {
			BeforeEach(func() {
				testServer = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
					defer GinkgoRecover()

					// Set success status code with no respoonse body
					res.WriteHeader(200)
				}))
			})
			It(`Invoke GetQueueManagerConnectionInfo successfully`, func() {
				mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1(&mqcloudv1.MqcloudV1Options{
					URL:            testServer.URL,
					Authenticator:  &core.NoAuthAuthenticator{},
					AcceptLanguage: core.StringPtr(acceptLanguage),
				})
				Expect(serviceErr).To(BeNil())
				Expect(mqcloudService).ToNot(BeNil())

				// Construct an instance of the GetQueueManagerConnectionInfoOptions model
				getQueueManagerConnectionInfoOptionsModel := new(mqcloudv1.GetQueueManagerConnectionInfoOptions)
				getQueueManagerConnectionInfoOptionsModel.ServiceInstanceGuid = core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				getQueueManagerConnectionInfoOptionsModel.QueueManagerID = core.StringPtr("b8e1aeda078009cf3db74e90d5d42328")
				getQueueManagerConnectionInfoOptionsModel.Headers = map[string]string{"x-custom-header": "x-custom-value"}

				// Invoke operation
				result, response, operationErr := mqcloudService.GetQueueManagerConnectionInfo(getQueueManagerConnectionInfoOptionsModel)
				Expect(operationErr).To(BeNil())
				Expect(response).ToNot(BeNil())

				// Verify a nil result
				Expect(result).To(BeNil())
			})
			AfterEach(func() {
				testServer.Close()
			})
		})
	})
	Describe(`GetQueueManagerStatus(getQueueManagerStatusOptions *GetQueueManagerStatusOptions) - Operation response error`, func() {
		acceptLanguage := "en-US,en;q=0.5"
		getQueueManagerStatusPath := "/v1/a2b4d4bc-dadb-4637-bcec-9b7d1e723af8/queue_managers/b8e1aeda078009cf3db74e90d5d42328/status"
		Context(`Using mock server endpoint with invalid JSON response`, func() {
			BeforeEach(func() {
				testServer = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
					defer GinkgoRecover()

					// Verify the contents of the request
					Expect(req.URL.EscapedPath()).To(Equal(getQueueManagerStatusPath))
					Expect(req.Method).To(Equal("GET"))
					Expect(req.Header["Accept-Language"]).ToNot(BeNil())
					Expect(req.Header["Accept-Language"][0]).To(Equal(fmt.Sprintf("%v", "en-US,en;q=0.5")))
					res.Header().Set("Content-type", "application/json")
					res.WriteHeader(200)
					fmt.Fprint(res, `} this is not valid json {`)
				}))
			})
			It(`Invoke GetQueueManagerStatus with error: Operation response processing error`, func() {
				mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1(&mqcloudv1.MqcloudV1Options{
					URL:            testServer.URL,
					Authenticator:  &core.NoAuthAuthenticator{},
					AcceptLanguage: core.StringPtr(acceptLanguage),
				})
				Expect(serviceErr).To(BeNil())
				Expect(mqcloudService).ToNot(BeNil())

				// Construct an instance of the GetQueueManagerStatusOptions model
				getQueueManagerStatusOptionsModel := new(mqcloudv1.GetQueueManagerStatusOptions)
				getQueueManagerStatusOptionsModel.ServiceInstanceGuid = core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				getQueueManagerStatusOptionsModel.QueueManagerID = core.StringPtr("b8e1aeda078009cf3db74e90d5d42328")
				getQueueManagerStatusOptionsModel.Headers = map[string]string{"x-custom-header": "x-custom-value"}
				// Expect response parsing to fail since we are receiving a text/plain response
				result, response, operationErr := mqcloudService.GetQueueManagerStatus(getQueueManagerStatusOptionsModel)
				Expect(operationErr).ToNot(BeNil())
				Expect(response).ToNot(BeNil())
				Expect(result).To(BeNil())

				// Enable retries and test again
				mqcloudService.EnableRetries(0, 0)
				result, response, operationErr = mqcloudService.GetQueueManagerStatus(getQueueManagerStatusOptionsModel)
				Expect(operationErr).ToNot(BeNil())
				Expect(response).ToNot(BeNil())
				Expect(result).To(BeNil())
			})
			AfterEach(func() {
				testServer.Close()
			})
		})
	})
	Describe(`GetQueueManagerStatus(getQueueManagerStatusOptions *GetQueueManagerStatusOptions)`, func() {
		acceptLanguage := "en-US,en;q=0.5"
		getQueueManagerStatusPath := "/v1/a2b4d4bc-dadb-4637-bcec-9b7d1e723af8/queue_managers/b8e1aeda078009cf3db74e90d5d42328/status"
		Context(`Using mock server endpoint with timeout`, func() {
			BeforeEach(func() {
				testServer = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
					defer GinkgoRecover()

					// Verify the contents of the request
					Expect(req.URL.EscapedPath()).To(Equal(getQueueManagerStatusPath))
					Expect(req.Method).To(Equal("GET"))

					Expect(req.Header["Accept-Language"]).ToNot(BeNil())
					Expect(req.Header["Accept-Language"][0]).To(Equal(fmt.Sprintf("%v", "en-US,en;q=0.5")))
					// Sleep a short time to support a timeout test
					time.Sleep(100 * time.Millisecond)

					// Set mock response
					res.Header().Set("Content-type", "application/json")
					res.WriteHeader(200)
					fmt.Fprintf(res, "%s", `{"status": "initializing"}`)
				}))
			})
			It(`Invoke GetQueueManagerStatus successfully with retries`, func() {
				mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1(&mqcloudv1.MqcloudV1Options{
					URL:            testServer.URL,
					Authenticator:  &core.NoAuthAuthenticator{},
					AcceptLanguage: core.StringPtr(acceptLanguage),
				})
				Expect(serviceErr).To(BeNil())
				Expect(mqcloudService).ToNot(BeNil())
				mqcloudService.EnableRetries(0, 0)

				// Construct an instance of the GetQueueManagerStatusOptions model
				getQueueManagerStatusOptionsModel := new(mqcloudv1.GetQueueManagerStatusOptions)
				getQueueManagerStatusOptionsModel.ServiceInstanceGuid = core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				getQueueManagerStatusOptionsModel.QueueManagerID = core.StringPtr("b8e1aeda078009cf3db74e90d5d42328")
				getQueueManagerStatusOptionsModel.Headers = map[string]string{"x-custom-header": "x-custom-value"}

				// Invoke operation with a Context to test a timeout error
				ctx, cancelFunc := context.WithTimeout(context.Background(), 80*time.Millisecond)
				defer cancelFunc()
				_, _, operationErr := mqcloudService.GetQueueManagerStatusWithContext(ctx, getQueueManagerStatusOptionsModel)
				Expect(operationErr).ToNot(BeNil())
				Expect(operationErr.Error()).To(ContainSubstring("deadline exceeded"))

				// Disable retries and test again
				mqcloudService.DisableRetries()
				result, response, operationErr := mqcloudService.GetQueueManagerStatus(getQueueManagerStatusOptionsModel)
				Expect(operationErr).To(BeNil())
				Expect(response).ToNot(BeNil())
				Expect(result).ToNot(BeNil())

				// Re-test the timeout error with retries disabled
				ctx, cancelFunc2 := context.WithTimeout(context.Background(), 80*time.Millisecond)
				defer cancelFunc2()
				_, _, operationErr = mqcloudService.GetQueueManagerStatusWithContext(ctx, getQueueManagerStatusOptionsModel)
				Expect(operationErr).ToNot(BeNil())
				Expect(operationErr.Error()).To(ContainSubstring("deadline exceeded"))
			})
			AfterEach(func() {
				testServer.Close()
			})
		})
		Context(`Using mock server endpoint`, func() {
			BeforeEach(func() {
				testServer = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
					defer GinkgoRecover()

					// Verify the contents of the request
					Expect(req.URL.EscapedPath()).To(Equal(getQueueManagerStatusPath))
					Expect(req.Method).To(Equal("GET"))

					Expect(req.Header["Accept-Language"]).ToNot(BeNil())
					Expect(req.Header["Accept-Language"][0]).To(Equal(fmt.Sprintf("%v", "en-US,en;q=0.5")))
					// Set mock response
					res.Header().Set("Content-type", "application/json")
					res.WriteHeader(200)
					fmt.Fprintf(res, "%s", `{"status": "initializing"}`)
				}))
			})
			It(`Invoke GetQueueManagerStatus successfully`, func() {
				mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1(&mqcloudv1.MqcloudV1Options{
					URL:            testServer.URL,
					Authenticator:  &core.NoAuthAuthenticator{},
					AcceptLanguage: core.StringPtr(acceptLanguage),
				})
				Expect(serviceErr).To(BeNil())
				Expect(mqcloudService).ToNot(BeNil())

				// Invoke operation with nil options model (negative test)
				result, response, operationErr := mqcloudService.GetQueueManagerStatus(nil)
				Expect(operationErr).NotTo(BeNil())
				Expect(response).To(BeNil())
				Expect(result).To(BeNil())

				// Construct an instance of the GetQueueManagerStatusOptions model
				getQueueManagerStatusOptionsModel := new(mqcloudv1.GetQueueManagerStatusOptions)
				getQueueManagerStatusOptionsModel.ServiceInstanceGuid = core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				getQueueManagerStatusOptionsModel.QueueManagerID = core.StringPtr("b8e1aeda078009cf3db74e90d5d42328")
				getQueueManagerStatusOptionsModel.Headers = map[string]string{"x-custom-header": "x-custom-value"}

				// Invoke operation with valid options model (positive test)
				result, response, operationErr = mqcloudService.GetQueueManagerStatus(getQueueManagerStatusOptionsModel)
				Expect(operationErr).To(BeNil())
				Expect(response).ToNot(BeNil())
				Expect(result).ToNot(BeNil())

			})
			It(`Invoke GetQueueManagerStatus with error: Operation validation and request error`, func() {
				mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1(&mqcloudv1.MqcloudV1Options{
					URL:            testServer.URL,
					Authenticator:  &core.NoAuthAuthenticator{},
					AcceptLanguage: core.StringPtr(acceptLanguage),
				})
				Expect(serviceErr).To(BeNil())
				Expect(mqcloudService).ToNot(BeNil())

				// Construct an instance of the GetQueueManagerStatusOptions model
				getQueueManagerStatusOptionsModel := new(mqcloudv1.GetQueueManagerStatusOptions)
				getQueueManagerStatusOptionsModel.ServiceInstanceGuid = core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				getQueueManagerStatusOptionsModel.QueueManagerID = core.StringPtr("b8e1aeda078009cf3db74e90d5d42328")
				getQueueManagerStatusOptionsModel.Headers = map[string]string{"x-custom-header": "x-custom-value"}
				// Invoke operation with empty URL (negative test)
				err := mqcloudService.SetServiceURL("")
				Expect(err).To(BeNil())
				result, response, operationErr := mqcloudService.GetQueueManagerStatus(getQueueManagerStatusOptionsModel)
				Expect(operationErr).ToNot(BeNil())
				Expect(operationErr.Error()).To(ContainSubstring(core.ERRORMSG_SERVICE_URL_MISSING))
				Expect(response).To(BeNil())
				Expect(result).To(BeNil())
				// Construct a second instance of the GetQueueManagerStatusOptions model with no property values
				getQueueManagerStatusOptionsModelNew := new(mqcloudv1.GetQueueManagerStatusOptions)
				// Invoke operation with invalid model (negative test)
				result, response, operationErr = mqcloudService.GetQueueManagerStatus(getQueueManagerStatusOptionsModelNew)
				Expect(operationErr).ToNot(BeNil())
				Expect(response).To(BeNil())
				Expect(result).To(BeNil())
			})
			AfterEach(func() {
				testServer.Close()
			})
		})
		Context(`Using mock server endpoint with missing response body`, func() {
			BeforeEach(func() {
				testServer = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
					defer GinkgoRecover()

					// Set success status code with no respoonse body
					res.WriteHeader(200)
				}))
			})
			It(`Invoke GetQueueManagerStatus successfully`, func() {
				mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1(&mqcloudv1.MqcloudV1Options{
					URL:            testServer.URL,
					Authenticator:  &core.NoAuthAuthenticator{},
					AcceptLanguage: core.StringPtr(acceptLanguage),
				})
				Expect(serviceErr).To(BeNil())
				Expect(mqcloudService).ToNot(BeNil())

				// Construct an instance of the GetQueueManagerStatusOptions model
				getQueueManagerStatusOptionsModel := new(mqcloudv1.GetQueueManagerStatusOptions)
				getQueueManagerStatusOptionsModel.ServiceInstanceGuid = core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				getQueueManagerStatusOptionsModel.QueueManagerID = core.StringPtr("b8e1aeda078009cf3db74e90d5d42328")
				getQueueManagerStatusOptionsModel.Headers = map[string]string{"x-custom-header": "x-custom-value"}

				// Invoke operation
				result, response, operationErr := mqcloudService.GetQueueManagerStatus(getQueueManagerStatusOptionsModel)
				Expect(operationErr).To(BeNil())
				Expect(response).ToNot(BeNil())

				// Verify a nil result
				Expect(result).To(BeNil())
			})
			AfterEach(func() {
				testServer.Close()
			})
		})
	})
	Describe(`ListUsers(listUsersOptions *ListUsersOptions) - Operation response error`, func() {
		acceptLanguage := "en-US,en;q=0.5"
		listUsersPath := "/v1/a2b4d4bc-dadb-4637-bcec-9b7d1e723af8/users"
		Context(`Using mock server endpoint with invalid JSON response`, func() {
			BeforeEach(func() {
				testServer = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
					defer GinkgoRecover()

					// Verify the contents of the request
					Expect(req.URL.EscapedPath()).To(Equal(listUsersPath))
					Expect(req.Method).To(Equal("GET"))
					Expect(req.Header["Accept-Language"]).ToNot(BeNil())
					Expect(req.Header["Accept-Language"][0]).To(Equal(fmt.Sprintf("%v", "en-US,en;q=0.5")))
					Expect(req.URL.Query()["offset"]).To(Equal([]string{fmt.Sprint(int64(0))}))
					Expect(req.URL.Query()["limit"]).To(Equal([]string{fmt.Sprint(int64(10))}))
					res.Header().Set("Content-type", "application/json")
					res.WriteHeader(200)
					fmt.Fprint(res, `} this is not valid json {`)
				}))
			})
			It(`Invoke ListUsers with error: Operation response processing error`, func() {
				mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1(&mqcloudv1.MqcloudV1Options{
					URL:            testServer.URL,
					Authenticator:  &core.NoAuthAuthenticator{},
					AcceptLanguage: core.StringPtr(acceptLanguage),
				})
				Expect(serviceErr).To(BeNil())
				Expect(mqcloudService).ToNot(BeNil())

				// Construct an instance of the ListUsersOptions model
				listUsersOptionsModel := new(mqcloudv1.ListUsersOptions)
				listUsersOptionsModel.ServiceInstanceGuid = core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				listUsersOptionsModel.Offset = core.Int64Ptr(int64(0))
				listUsersOptionsModel.Limit = core.Int64Ptr(int64(10))
				listUsersOptionsModel.Headers = map[string]string{"x-custom-header": "x-custom-value"}
				// Expect response parsing to fail since we are receiving a text/plain response
				result, response, operationErr := mqcloudService.ListUsers(listUsersOptionsModel)
				Expect(operationErr).ToNot(BeNil())
				Expect(response).ToNot(BeNil())
				Expect(result).To(BeNil())

				// Enable retries and test again
				mqcloudService.EnableRetries(0, 0)
				result, response, operationErr = mqcloudService.ListUsers(listUsersOptionsModel)
				Expect(operationErr).ToNot(BeNil())
				Expect(response).ToNot(BeNil())
				Expect(result).To(BeNil())
			})
			AfterEach(func() {
				testServer.Close()
			})
		})
	})
	Describe(`ListUsers(listUsersOptions *ListUsersOptions)`, func() {
		acceptLanguage := "en-US,en;q=0.5"
		listUsersPath := "/v1/a2b4d4bc-dadb-4637-bcec-9b7d1e723af8/users"
		Context(`Using mock server endpoint with timeout`, func() {
			BeforeEach(func() {
				testServer = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
					defer GinkgoRecover()

					// Verify the contents of the request
					Expect(req.URL.EscapedPath()).To(Equal(listUsersPath))
					Expect(req.Method).To(Equal("GET"))

					Expect(req.Header["Accept-Language"]).ToNot(BeNil())
					Expect(req.Header["Accept-Language"][0]).To(Equal(fmt.Sprintf("%v", "en-US,en;q=0.5")))
					Expect(req.URL.Query()["offset"]).To(Equal([]string{fmt.Sprint(int64(0))}))
					Expect(req.URL.Query()["limit"]).To(Equal([]string{fmt.Sprint(int64(10))}))
					// Sleep a short time to support a timeout test
					time.Sleep(100 * time.Millisecond)

					// Set mock response
					res.Header().Set("Content-type", "application/json")
					res.WriteHeader(200)
					fmt.Fprintf(res, "%s", `{"offset": 6, "limit": 5, "first": {"href": "Href"}, "next": {"href": "Href"}, "previous": {"href": "Href"}, "users": [{"id": "ID", "name": "Name", "email": "user@host.org", "iam_service_id": "IamServiceID", "roles": ["Roles"], "iam_managed": true, "href": "Href"}]}`)
				}))
			})
			It(`Invoke ListUsers successfully with retries`, func() {
				mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1(&mqcloudv1.MqcloudV1Options{
					URL:            testServer.URL,
					Authenticator:  &core.NoAuthAuthenticator{},
					AcceptLanguage: core.StringPtr(acceptLanguage),
				})
				Expect(serviceErr).To(BeNil())
				Expect(mqcloudService).ToNot(BeNil())
				mqcloudService.EnableRetries(0, 0)

				// Construct an instance of the ListUsersOptions model
				listUsersOptionsModel := new(mqcloudv1.ListUsersOptions)
				listUsersOptionsModel.ServiceInstanceGuid = core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				listUsersOptionsModel.Offset = core.Int64Ptr(int64(0))
				listUsersOptionsModel.Limit = core.Int64Ptr(int64(10))
				listUsersOptionsModel.Headers = map[string]string{"x-custom-header": "x-custom-value"}

				// Invoke operation with a Context to test a timeout error
				ctx, cancelFunc := context.WithTimeout(context.Background(), 80*time.Millisecond)
				defer cancelFunc()
				_, _, operationErr := mqcloudService.ListUsersWithContext(ctx, listUsersOptionsModel)
				Expect(operationErr).ToNot(BeNil())
				Expect(operationErr.Error()).To(ContainSubstring("deadline exceeded"))

				// Disable retries and test again
				mqcloudService.DisableRetries()
				result, response, operationErr := mqcloudService.ListUsers(listUsersOptionsModel)
				Expect(operationErr).To(BeNil())
				Expect(response).ToNot(BeNil())
				Expect(result).ToNot(BeNil())

				// Re-test the timeout error with retries disabled
				ctx, cancelFunc2 := context.WithTimeout(context.Background(), 80*time.Millisecond)
				defer cancelFunc2()
				_, _, operationErr = mqcloudService.ListUsersWithContext(ctx, listUsersOptionsModel)
				Expect(operationErr).ToNot(BeNil())
				Expect(operationErr.Error()).To(ContainSubstring("deadline exceeded"))
			})
			AfterEach(func() {
				testServer.Close()
			})
		})
		Context(`Using mock server endpoint`, func() {
			BeforeEach(func() {
				testServer = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
					defer GinkgoRecover()

					// Verify the contents of the request
					Expect(req.URL.EscapedPath()).To(Equal(listUsersPath))
					Expect(req.Method).To(Equal("GET"))

					Expect(req.Header["Accept-Language"]).ToNot(BeNil())
					Expect(req.Header["Accept-Language"][0]).To(Equal(fmt.Sprintf("%v", "en-US,en;q=0.5")))
					Expect(req.URL.Query()["offset"]).To(Equal([]string{fmt.Sprint(int64(0))}))
					Expect(req.URL.Query()["limit"]).To(Equal([]string{fmt.Sprint(int64(10))}))
					// Set mock response
					res.Header().Set("Content-type", "application/json")
					res.WriteHeader(200)
					fmt.Fprintf(res, "%s", `{"offset": 6, "limit": 5, "first": {"href": "Href"}, "next": {"href": "Href"}, "previous": {"href": "Href"}, "users": [{"id": "ID", "name": "Name", "email": "user@host.org", "iam_service_id": "IamServiceID", "roles": ["Roles"], "iam_managed": true, "href": "Href"}]}`)
				}))
			})
			It(`Invoke ListUsers successfully`, func() {
				mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1(&mqcloudv1.MqcloudV1Options{
					URL:            testServer.URL,
					Authenticator:  &core.NoAuthAuthenticator{},
					AcceptLanguage: core.StringPtr(acceptLanguage),
				})
				Expect(serviceErr).To(BeNil())
				Expect(mqcloudService).ToNot(BeNil())

				// Invoke operation with nil options model (negative test)
				result, response, operationErr := mqcloudService.ListUsers(nil)
				Expect(operationErr).NotTo(BeNil())
				Expect(response).To(BeNil())
				Expect(result).To(BeNil())

				// Construct an instance of the ListUsersOptions model
				listUsersOptionsModel := new(mqcloudv1.ListUsersOptions)
				listUsersOptionsModel.ServiceInstanceGuid = core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				listUsersOptionsModel.Offset = core.Int64Ptr(int64(0))
				listUsersOptionsModel.Limit = core.Int64Ptr(int64(10))
				listUsersOptionsModel.Headers = map[string]string{"x-custom-header": "x-custom-value"}

				// Invoke operation with valid options model (positive test)
				result, response, operationErr = mqcloudService.ListUsers(listUsersOptionsModel)
				Expect(operationErr).To(BeNil())
				Expect(response).ToNot(BeNil())
				Expect(result).ToNot(BeNil())

			})
			It(`Invoke ListUsers with error: Operation validation and request error`, func() {
				mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1(&mqcloudv1.MqcloudV1Options{
					URL:            testServer.URL,
					Authenticator:  &core.NoAuthAuthenticator{},
					AcceptLanguage: core.StringPtr(acceptLanguage),
				})
				Expect(serviceErr).To(BeNil())
				Expect(mqcloudService).ToNot(BeNil())

				// Construct an instance of the ListUsersOptions model
				listUsersOptionsModel := new(mqcloudv1.ListUsersOptions)
				listUsersOptionsModel.ServiceInstanceGuid = core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				listUsersOptionsModel.Offset = core.Int64Ptr(int64(0))
				listUsersOptionsModel.Limit = core.Int64Ptr(int64(10))
				listUsersOptionsModel.Headers = map[string]string{"x-custom-header": "x-custom-value"}
				// Invoke operation with empty URL (negative test)
				err := mqcloudService.SetServiceURL("")
				Expect(err).To(BeNil())
				result, response, operationErr := mqcloudService.ListUsers(listUsersOptionsModel)
				Expect(operationErr).ToNot(BeNil())
				Expect(operationErr.Error()).To(ContainSubstring(core.ERRORMSG_SERVICE_URL_MISSING))
				Expect(response).To(BeNil())
				Expect(result).To(BeNil())
				// Construct a second instance of the ListUsersOptions model with no property values
				listUsersOptionsModelNew := new(mqcloudv1.ListUsersOptions)
				// Invoke operation with invalid model (negative test)
				result, response, operationErr = mqcloudService.ListUsers(listUsersOptionsModelNew)
				Expect(operationErr).ToNot(BeNil())
				Expect(response).To(BeNil())
				Expect(result).To(BeNil())
			})
			AfterEach(func() {
				testServer.Close()
			})
		})
		Context(`Using mock server endpoint with missing response body`, func() {
			BeforeEach(func() {
				testServer = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
					defer GinkgoRecover()

					// Set success status code with no respoonse body
					res.WriteHeader(200)
				}))
			})
			It(`Invoke ListUsers successfully`, func() {
				mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1(&mqcloudv1.MqcloudV1Options{
					URL:            testServer.URL,
					Authenticator:  &core.NoAuthAuthenticator{},
					AcceptLanguage: core.StringPtr(acceptLanguage),
				})
				Expect(serviceErr).To(BeNil())
				Expect(mqcloudService).ToNot(BeNil())

				// Construct an instance of the ListUsersOptions model
				listUsersOptionsModel := new(mqcloudv1.ListUsersOptions)
				listUsersOptionsModel.ServiceInstanceGuid = core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				listUsersOptionsModel.Offset = core.Int64Ptr(int64(0))
				listUsersOptionsModel.Limit = core.Int64Ptr(int64(10))
				listUsersOptionsModel.Headers = map[string]string{"x-custom-header": "x-custom-value"}

				// Invoke operation
				result, response, operationErr := mqcloudService.ListUsers(listUsersOptionsModel)
				Expect(operationErr).To(BeNil())
				Expect(response).ToNot(BeNil())

				// Verify a nil result
				Expect(result).To(BeNil())
			})
			AfterEach(func() {
				testServer.Close()
			})
		})
		Context(`Test pagination helper method on response`, func() {
			It(`Invoke GetNextOffset successfully`, func() {
				responseObject := new(mqcloudv1.UserDetailsCollection)
				nextObject := new(mqcloudv1.Next)
				nextObject.Href = core.StringPtr("ibm.com?offset=135")
				responseObject.Next = nextObject

				value, err := responseObject.GetNextOffset()
				Expect(err).To(BeNil())
				Expect(value).To(Equal(core.Int64Ptr(int64(135))))
			})
			It(`Invoke GetNextOffset without a "Next" property in the response`, func() {
				responseObject := new(mqcloudv1.UserDetailsCollection)

				value, err := responseObject.GetNextOffset()
				Expect(err).To(BeNil())
				Expect(value).To(BeNil())
			})
			It(`Invoke GetNextOffset without any query params in the "Next" URL`, func() {
				responseObject := new(mqcloudv1.UserDetailsCollection)
				nextObject := new(mqcloudv1.Next)
				nextObject.Href = core.StringPtr("ibm.com")
				responseObject.Next = nextObject

				value, err := responseObject.GetNextOffset()
				Expect(err).To(BeNil())
				Expect(value).To(BeNil())
			})
			It(`Invoke GetNextOffset with a non-integer query param in the "Next" URL`, func() {
				responseObject := new(mqcloudv1.UserDetailsCollection)
				nextObject := new(mqcloudv1.Next)
				nextObject.Href = core.StringPtr("ibm.com?offset=tiger")
				responseObject.Next = nextObject

				value, err := responseObject.GetNextOffset()
				Expect(err).NotTo(BeNil())
				Expect(value).To(BeNil())
			})
		})
		Context(`Using mock server endpoint - paginated response`, func() {
			BeforeEach(func() {
				var requestNumber int = 0
				testServer = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
					defer GinkgoRecover()

					// Verify the contents of the request
					Expect(req.URL.EscapedPath()).To(Equal(listUsersPath))
					Expect(req.Method).To(Equal("GET"))

					// Set mock response
					res.Header().Set("Content-type", "application/json")
					res.WriteHeader(200)
					requestNumber++
					if requestNumber == 1 {
						fmt.Fprintf(res, "%s", `{"next":{"href":"https://myhost.com/somePath?offset=1"},"total_count":2,"limit":1,"users":[{"id":"ID","name":"Name","email":"user@host.org","iam_service_id":"IamServiceID","roles":["Roles"],"iam_managed":true,"href":"Href"}]}`)
					} else if requestNumber == 2 {
						fmt.Fprintf(res, "%s", `{"total_count":2,"limit":1,"users":[{"id":"ID","name":"Name","email":"user@host.org","iam_service_id":"IamServiceID","roles":["Roles"],"iam_managed":true,"href":"Href"}]}`)
					} else {
						res.WriteHeader(400)
					}
				}))
			})
			It(`Use UsersPager.GetNext successfully`, func() {
				mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1(&mqcloudv1.MqcloudV1Options{
					URL:            testServer.URL,
					Authenticator:  &core.NoAuthAuthenticator{},
					AcceptLanguage: core.StringPtr(acceptLanguage),
				})
				Expect(serviceErr).To(BeNil())
				Expect(mqcloudService).ToNot(BeNil())

				listUsersOptionsModel := &mqcloudv1.ListUsersOptions{
					ServiceInstanceGuid: core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8"),
					Limit:               core.Int64Ptr(int64(10)),
				}

				pager, err := mqcloudService.NewUsersPager(listUsersOptionsModel)
				Expect(err).To(BeNil())
				Expect(pager).ToNot(BeNil())

				var allResults []mqcloudv1.UserDetails
				for pager.HasNext() {
					nextPage, err := pager.GetNext()
					Expect(err).To(BeNil())
					Expect(nextPage).ToNot(BeNil())
					allResults = append(allResults, nextPage...)
				}
				Expect(len(allResults)).To(Equal(2))
			})
			It(`Use UsersPager.GetAll successfully`, func() {
				mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1(&mqcloudv1.MqcloudV1Options{
					URL:            testServer.URL,
					Authenticator:  &core.NoAuthAuthenticator{},
					AcceptLanguage: core.StringPtr(acceptLanguage),
				})
				Expect(serviceErr).To(BeNil())
				Expect(mqcloudService).ToNot(BeNil())

				listUsersOptionsModel := &mqcloudv1.ListUsersOptions{
					ServiceInstanceGuid: core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8"),
					Limit:               core.Int64Ptr(int64(10)),
				}

				pager, err := mqcloudService.NewUsersPager(listUsersOptionsModel)
				Expect(err).To(BeNil())
				Expect(pager).ToNot(BeNil())

				allResults, err := pager.GetAll()
				Expect(err).To(BeNil())
				Expect(allResults).ToNot(BeNil())
				Expect(len(allResults)).To(Equal(2))
			})
		})
	})
	Describe(`CreateUser(createUserOptions *CreateUserOptions) - Operation response error`, func() {
		acceptLanguage := "en-US,en;q=0.5"
		createUserPath := "/v1/a2b4d4bc-dadb-4637-bcec-9b7d1e723af8/users"
		Context(`Using mock server endpoint with invalid JSON response`, func() {
			BeforeEach(func() {
				testServer = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
					defer GinkgoRecover()

					// Verify the contents of the request
					Expect(req.URL.EscapedPath()).To(Equal(createUserPath))
					Expect(req.Method).To(Equal("POST"))
					Expect(req.Header["Accept-Language"]).ToNot(BeNil())
					Expect(req.Header["Accept-Language"][0]).To(Equal(fmt.Sprintf("%v", "en-US,en;q=0.5")))
					res.Header().Set("Content-type", "application/json")
					res.WriteHeader(201)
					fmt.Fprint(res, `} this is not valid json {`)
				}))
			})
			It(`Invoke CreateUser with error: Operation response processing error`, func() {
				mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1(&mqcloudv1.MqcloudV1Options{
					URL:            testServer.URL,
					Authenticator:  &core.NoAuthAuthenticator{},
					AcceptLanguage: core.StringPtr(acceptLanguage),
				})
				Expect(serviceErr).To(BeNil())
				Expect(mqcloudService).ToNot(BeNil())

				// Construct an instance of the CreateUserOptions model
				createUserOptionsModel := new(mqcloudv1.CreateUserOptions)
				createUserOptionsModel.ServiceInstanceGuid = core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				createUserOptionsModel.Email = core.StringPtr("testuser@ibm.com")
				createUserOptionsModel.Name = core.StringPtr("testuser")
				createUserOptionsModel.Headers = map[string]string{"x-custom-header": "x-custom-value"}
				// Expect response parsing to fail since we are receiving a text/plain response
				result, response, operationErr := mqcloudService.CreateUser(createUserOptionsModel)
				Expect(operationErr).ToNot(BeNil())
				Expect(response).ToNot(BeNil())
				Expect(result).To(BeNil())

				// Enable retries and test again
				mqcloudService.EnableRetries(0, 0)
				result, response, operationErr = mqcloudService.CreateUser(createUserOptionsModel)
				Expect(operationErr).ToNot(BeNil())
				Expect(response).ToNot(BeNil())
				Expect(result).To(BeNil())
			})
			AfterEach(func() {
				testServer.Close()
			})
		})
	})
	Describe(`CreateUser(createUserOptions *CreateUserOptions)`, func() {
		acceptLanguage := "en-US,en;q=0.5"
		createUserPath := "/v1/a2b4d4bc-dadb-4637-bcec-9b7d1e723af8/users"
		Context(`Using mock server endpoint with timeout`, func() {
			BeforeEach(func() {
				testServer = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
					defer GinkgoRecover()

					// Verify the contents of the request
					Expect(req.URL.EscapedPath()).To(Equal(createUserPath))
					Expect(req.Method).To(Equal("POST"))

					// For gzip-disabled operation, verify Content-Encoding is not set.
					Expect(req.Header.Get("Content-Encoding")).To(BeEmpty())

					// If there is a body, then make sure we can read it
					bodyBuf := new(bytes.Buffer)
					if req.Header.Get("Content-Encoding") == "gzip" {
						body, err := core.NewGzipDecompressionReader(req.Body)
						Expect(err).To(BeNil())
						_, err = bodyBuf.ReadFrom(body)
						Expect(err).To(BeNil())
					} else {
						_, err := bodyBuf.ReadFrom(req.Body)
						Expect(err).To(BeNil())
					}
					fmt.Fprintf(GinkgoWriter, "  Request body: %s", bodyBuf.String())

					Expect(req.Header["Accept-Language"]).ToNot(BeNil())
					Expect(req.Header["Accept-Language"][0]).To(Equal(fmt.Sprintf("%v", "en-US,en;q=0.5")))
					// Sleep a short time to support a timeout test
					time.Sleep(100 * time.Millisecond)

					// Set mock response
					res.Header().Set("Content-type", "application/json")
					res.WriteHeader(201)
					fmt.Fprintf(res, "%s", `{"id": "ID", "name": "Name", "email": "user@host.org", "iam_service_id": "IamServiceID", "roles": ["Roles"], "iam_managed": true, "href": "Href"}`)
				}))
			})
			It(`Invoke CreateUser successfully with retries`, func() {
				mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1(&mqcloudv1.MqcloudV1Options{
					URL:            testServer.URL,
					Authenticator:  &core.NoAuthAuthenticator{},
					AcceptLanguage: core.StringPtr(acceptLanguage),
				})
				Expect(serviceErr).To(BeNil())
				Expect(mqcloudService).ToNot(BeNil())
				mqcloudService.EnableRetries(0, 0)

				// Construct an instance of the CreateUserOptions model
				createUserOptionsModel := new(mqcloudv1.CreateUserOptions)
				createUserOptionsModel.ServiceInstanceGuid = core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				createUserOptionsModel.Email = core.StringPtr("testuser@ibm.com")
				createUserOptionsModel.Name = core.StringPtr("testuser")
				createUserOptionsModel.Headers = map[string]string{"x-custom-header": "x-custom-value"}

				// Invoke operation with a Context to test a timeout error
				ctx, cancelFunc := context.WithTimeout(context.Background(), 80*time.Millisecond)
				defer cancelFunc()
				_, _, operationErr := mqcloudService.CreateUserWithContext(ctx, createUserOptionsModel)
				Expect(operationErr).ToNot(BeNil())
				Expect(operationErr.Error()).To(ContainSubstring("deadline exceeded"))

				// Disable retries and test again
				mqcloudService.DisableRetries()
				result, response, operationErr := mqcloudService.CreateUser(createUserOptionsModel)
				Expect(operationErr).To(BeNil())
				Expect(response).ToNot(BeNil())
				Expect(result).ToNot(BeNil())

				// Re-test the timeout error with retries disabled
				ctx, cancelFunc2 := context.WithTimeout(context.Background(), 80*time.Millisecond)
				defer cancelFunc2()
				_, _, operationErr = mqcloudService.CreateUserWithContext(ctx, createUserOptionsModel)
				Expect(operationErr).ToNot(BeNil())
				Expect(operationErr.Error()).To(ContainSubstring("deadline exceeded"))
			})
			AfterEach(func() {
				testServer.Close()
			})
		})
		Context(`Using mock server endpoint`, func() {
			BeforeEach(func() {
				testServer = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
					defer GinkgoRecover()

					// Verify the contents of the request
					Expect(req.URL.EscapedPath()).To(Equal(createUserPath))
					Expect(req.Method).To(Equal("POST"))

					// For gzip-disabled operation, verify Content-Encoding is not set.
					Expect(req.Header.Get("Content-Encoding")).To(BeEmpty())

					// If there is a body, then make sure we can read it
					bodyBuf := new(bytes.Buffer)
					if req.Header.Get("Content-Encoding") == "gzip" {
						body, err := core.NewGzipDecompressionReader(req.Body)
						Expect(err).To(BeNil())
						_, err = bodyBuf.ReadFrom(body)
						Expect(err).To(BeNil())
					} else {
						_, err := bodyBuf.ReadFrom(req.Body)
						Expect(err).To(BeNil())
					}
					fmt.Fprintf(GinkgoWriter, "  Request body: %s", bodyBuf.String())

					Expect(req.Header["Accept-Language"]).ToNot(BeNil())
					Expect(req.Header["Accept-Language"][0]).To(Equal(fmt.Sprintf("%v", "en-US,en;q=0.5")))
					// Set mock response
					res.Header().Set("Content-type", "application/json")
					res.WriteHeader(201)
					fmt.Fprintf(res, "%s", `{"id": "ID", "name": "Name", "email": "user@host.org", "iam_service_id": "IamServiceID", "roles": ["Roles"], "iam_managed": true, "href": "Href"}`)
				}))
			})
			It(`Invoke CreateUser successfully`, func() {
				mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1(&mqcloudv1.MqcloudV1Options{
					URL:            testServer.URL,
					Authenticator:  &core.NoAuthAuthenticator{},
					AcceptLanguage: core.StringPtr(acceptLanguage),
				})
				Expect(serviceErr).To(BeNil())
				Expect(mqcloudService).ToNot(BeNil())

				// Invoke operation with nil options model (negative test)
				result, response, operationErr := mqcloudService.CreateUser(nil)
				Expect(operationErr).NotTo(BeNil())
				Expect(response).To(BeNil())
				Expect(result).To(BeNil())

				// Construct an instance of the CreateUserOptions model
				createUserOptionsModel := new(mqcloudv1.CreateUserOptions)
				createUserOptionsModel.ServiceInstanceGuid = core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				createUserOptionsModel.Email = core.StringPtr("testuser@ibm.com")
				createUserOptionsModel.Name = core.StringPtr("testuser")
				createUserOptionsModel.Headers = map[string]string{"x-custom-header": "x-custom-value"}

				// Invoke operation with valid options model (positive test)
				result, response, operationErr = mqcloudService.CreateUser(createUserOptionsModel)
				Expect(operationErr).To(BeNil())
				Expect(response).ToNot(BeNil())
				Expect(result).ToNot(BeNil())

			})
			It(`Invoke CreateUser with error: Operation validation and request error`, func() {
				mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1(&mqcloudv1.MqcloudV1Options{
					URL:            testServer.URL,
					Authenticator:  &core.NoAuthAuthenticator{},
					AcceptLanguage: core.StringPtr(acceptLanguage),
				})
				Expect(serviceErr).To(BeNil())
				Expect(mqcloudService).ToNot(BeNil())

				// Construct an instance of the CreateUserOptions model
				createUserOptionsModel := new(mqcloudv1.CreateUserOptions)
				createUserOptionsModel.ServiceInstanceGuid = core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				createUserOptionsModel.Email = core.StringPtr("testuser@ibm.com")
				createUserOptionsModel.Name = core.StringPtr("testuser")
				createUserOptionsModel.Headers = map[string]string{"x-custom-header": "x-custom-value"}
				// Invoke operation with empty URL (negative test)
				err := mqcloudService.SetServiceURL("")
				Expect(err).To(BeNil())
				result, response, operationErr := mqcloudService.CreateUser(createUserOptionsModel)
				Expect(operationErr).ToNot(BeNil())
				Expect(operationErr.Error()).To(ContainSubstring(core.ERRORMSG_SERVICE_URL_MISSING))
				Expect(response).To(BeNil())
				Expect(result).To(BeNil())
				// Construct a second instance of the CreateUserOptions model with no property values
				createUserOptionsModelNew := new(mqcloudv1.CreateUserOptions)
				// Invoke operation with invalid model (negative test)
				result, response, operationErr = mqcloudService.CreateUser(createUserOptionsModelNew)
				Expect(operationErr).ToNot(BeNil())
				Expect(response).To(BeNil())
				Expect(result).To(BeNil())
			})
			AfterEach(func() {
				testServer.Close()
			})
		})
		Context(`Using mock server endpoint with missing response body`, func() {
			BeforeEach(func() {
				testServer = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
					defer GinkgoRecover()

					// Set success status code with no respoonse body
					res.WriteHeader(201)
				}))
			})
			It(`Invoke CreateUser successfully`, func() {
				mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1(&mqcloudv1.MqcloudV1Options{
					URL:            testServer.URL,
					Authenticator:  &core.NoAuthAuthenticator{},
					AcceptLanguage: core.StringPtr(acceptLanguage),
				})
				Expect(serviceErr).To(BeNil())
				Expect(mqcloudService).ToNot(BeNil())

				// Construct an instance of the CreateUserOptions model
				createUserOptionsModel := new(mqcloudv1.CreateUserOptions)
				createUserOptionsModel.ServiceInstanceGuid = core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				createUserOptionsModel.Email = core.StringPtr("testuser@ibm.com")
				createUserOptionsModel.Name = core.StringPtr("testuser")
				createUserOptionsModel.Headers = map[string]string{"x-custom-header": "x-custom-value"}

				// Invoke operation
				result, response, operationErr := mqcloudService.CreateUser(createUserOptionsModel)
				Expect(operationErr).To(BeNil())
				Expect(response).ToNot(BeNil())

				// Verify a nil result
				Expect(result).To(BeNil())
			})
			AfterEach(func() {
				testServer.Close()
			})
		})
	})
	Describe(`GetUser(getUserOptions *GetUserOptions) - Operation response error`, func() {
		acceptLanguage := "en-US,en;q=0.5"
		getUserPath := "/v1/a2b4d4bc-dadb-4637-bcec-9b7d1e723af8/users/31a413dd84346effc8895b6ba4641641"
		Context(`Using mock server endpoint with invalid JSON response`, func() {
			BeforeEach(func() {
				testServer = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
					defer GinkgoRecover()

					// Verify the contents of the request
					Expect(req.URL.EscapedPath()).To(Equal(getUserPath))
					Expect(req.Method).To(Equal("GET"))
					Expect(req.Header["Accept-Language"]).ToNot(BeNil())
					Expect(req.Header["Accept-Language"][0]).To(Equal(fmt.Sprintf("%v", "en-US,en;q=0.5")))
					res.Header().Set("Content-type", "application/json")
					res.WriteHeader(200)
					fmt.Fprint(res, `} this is not valid json {`)
				}))
			})
			It(`Invoke GetUser with error: Operation response processing error`, func() {
				mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1(&mqcloudv1.MqcloudV1Options{
					URL:            testServer.URL,
					Authenticator:  &core.NoAuthAuthenticator{},
					AcceptLanguage: core.StringPtr(acceptLanguage),
				})
				Expect(serviceErr).To(BeNil())
				Expect(mqcloudService).ToNot(BeNil())

				// Construct an instance of the GetUserOptions model
				getUserOptionsModel := new(mqcloudv1.GetUserOptions)
				getUserOptionsModel.ServiceInstanceGuid = core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				getUserOptionsModel.UserID = core.StringPtr("31a413dd84346effc8895b6ba4641641")
				getUserOptionsModel.Headers = map[string]string{"x-custom-header": "x-custom-value"}
				// Expect response parsing to fail since we are receiving a text/plain response
				result, response, operationErr := mqcloudService.GetUser(getUserOptionsModel)
				Expect(operationErr).ToNot(BeNil())
				Expect(response).ToNot(BeNil())
				Expect(result).To(BeNil())

				// Enable retries and test again
				mqcloudService.EnableRetries(0, 0)
				result, response, operationErr = mqcloudService.GetUser(getUserOptionsModel)
				Expect(operationErr).ToNot(BeNil())
				Expect(response).ToNot(BeNil())
				Expect(result).To(BeNil())
			})
			AfterEach(func() {
				testServer.Close()
			})
		})
	})
	Describe(`GetUser(getUserOptions *GetUserOptions)`, func() {
		acceptLanguage := "en-US,en;q=0.5"
		getUserPath := "/v1/a2b4d4bc-dadb-4637-bcec-9b7d1e723af8/users/31a413dd84346effc8895b6ba4641641"
		Context(`Using mock server endpoint with timeout`, func() {
			BeforeEach(func() {
				testServer = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
					defer GinkgoRecover()

					// Verify the contents of the request
					Expect(req.URL.EscapedPath()).To(Equal(getUserPath))
					Expect(req.Method).To(Equal("GET"))

					Expect(req.Header["Accept-Language"]).ToNot(BeNil())
					Expect(req.Header["Accept-Language"][0]).To(Equal(fmt.Sprintf("%v", "en-US,en;q=0.5")))
					// Sleep a short time to support a timeout test
					time.Sleep(100 * time.Millisecond)

					// Set mock response
					res.Header().Set("Content-type", "application/json")
					res.WriteHeader(200)
					fmt.Fprintf(res, "%s", `{"id": "ID", "name": "Name", "email": "user@host.org", "iam_service_id": "IamServiceID", "roles": ["Roles"], "iam_managed": true, "href": "Href"}`)
				}))
			})
			It(`Invoke GetUser successfully with retries`, func() {
				mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1(&mqcloudv1.MqcloudV1Options{
					URL:            testServer.URL,
					Authenticator:  &core.NoAuthAuthenticator{},
					AcceptLanguage: core.StringPtr(acceptLanguage),
				})
				Expect(serviceErr).To(BeNil())
				Expect(mqcloudService).ToNot(BeNil())
				mqcloudService.EnableRetries(0, 0)

				// Construct an instance of the GetUserOptions model
				getUserOptionsModel := new(mqcloudv1.GetUserOptions)
				getUserOptionsModel.ServiceInstanceGuid = core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				getUserOptionsModel.UserID = core.StringPtr("31a413dd84346effc8895b6ba4641641")
				getUserOptionsModel.Headers = map[string]string{"x-custom-header": "x-custom-value"}

				// Invoke operation with a Context to test a timeout error
				ctx, cancelFunc := context.WithTimeout(context.Background(), 80*time.Millisecond)
				defer cancelFunc()
				_, _, operationErr := mqcloudService.GetUserWithContext(ctx, getUserOptionsModel)
				Expect(operationErr).ToNot(BeNil())
				Expect(operationErr.Error()).To(ContainSubstring("deadline exceeded"))

				// Disable retries and test again
				mqcloudService.DisableRetries()
				result, response, operationErr := mqcloudService.GetUser(getUserOptionsModel)
				Expect(operationErr).To(BeNil())
				Expect(response).ToNot(BeNil())
				Expect(result).ToNot(BeNil())

				// Re-test the timeout error with retries disabled
				ctx, cancelFunc2 := context.WithTimeout(context.Background(), 80*time.Millisecond)
				defer cancelFunc2()
				_, _, operationErr = mqcloudService.GetUserWithContext(ctx, getUserOptionsModel)
				Expect(operationErr).ToNot(BeNil())
				Expect(operationErr.Error()).To(ContainSubstring("deadline exceeded"))
			})
			AfterEach(func() {
				testServer.Close()
			})
		})
		Context(`Using mock server endpoint`, func() {
			BeforeEach(func() {
				testServer = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
					defer GinkgoRecover()

					// Verify the contents of the request
					Expect(req.URL.EscapedPath()).To(Equal(getUserPath))
					Expect(req.Method).To(Equal("GET"))

					Expect(req.Header["Accept-Language"]).ToNot(BeNil())
					Expect(req.Header["Accept-Language"][0]).To(Equal(fmt.Sprintf("%v", "en-US,en;q=0.5")))
					// Set mock response
					res.Header().Set("Content-type", "application/json")
					res.WriteHeader(200)
					fmt.Fprintf(res, "%s", `{"id": "ID", "name": "Name", "email": "user@host.org", "iam_service_id": "IamServiceID", "roles": ["Roles"], "iam_managed": true, "href": "Href"}`)
				}))
			})
			It(`Invoke GetUser successfully`, func() {
				mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1(&mqcloudv1.MqcloudV1Options{
					URL:            testServer.URL,
					Authenticator:  &core.NoAuthAuthenticator{},
					AcceptLanguage: core.StringPtr(acceptLanguage),
				})
				Expect(serviceErr).To(BeNil())
				Expect(mqcloudService).ToNot(BeNil())

				// Invoke operation with nil options model (negative test)
				result, response, operationErr := mqcloudService.GetUser(nil)
				Expect(operationErr).NotTo(BeNil())
				Expect(response).To(BeNil())
				Expect(result).To(BeNil())

				// Construct an instance of the GetUserOptions model
				getUserOptionsModel := new(mqcloudv1.GetUserOptions)
				getUserOptionsModel.ServiceInstanceGuid = core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				getUserOptionsModel.UserID = core.StringPtr("31a413dd84346effc8895b6ba4641641")
				getUserOptionsModel.Headers = map[string]string{"x-custom-header": "x-custom-value"}

				// Invoke operation with valid options model (positive test)
				result, response, operationErr = mqcloudService.GetUser(getUserOptionsModel)
				Expect(operationErr).To(BeNil())
				Expect(response).ToNot(BeNil())
				Expect(result).ToNot(BeNil())

			})
			It(`Invoke GetUser with error: Operation validation and request error`, func() {
				mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1(&mqcloudv1.MqcloudV1Options{
					URL:            testServer.URL,
					Authenticator:  &core.NoAuthAuthenticator{},
					AcceptLanguage: core.StringPtr(acceptLanguage),
				})
				Expect(serviceErr).To(BeNil())
				Expect(mqcloudService).ToNot(BeNil())

				// Construct an instance of the GetUserOptions model
				getUserOptionsModel := new(mqcloudv1.GetUserOptions)
				getUserOptionsModel.ServiceInstanceGuid = core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				getUserOptionsModel.UserID = core.StringPtr("31a413dd84346effc8895b6ba4641641")
				getUserOptionsModel.Headers = map[string]string{"x-custom-header": "x-custom-value"}
				// Invoke operation with empty URL (negative test)
				err := mqcloudService.SetServiceURL("")
				Expect(err).To(BeNil())
				result, response, operationErr := mqcloudService.GetUser(getUserOptionsModel)
				Expect(operationErr).ToNot(BeNil())
				Expect(operationErr.Error()).To(ContainSubstring(core.ERRORMSG_SERVICE_URL_MISSING))
				Expect(response).To(BeNil())
				Expect(result).To(BeNil())
				// Construct a second instance of the GetUserOptions model with no property values
				getUserOptionsModelNew := new(mqcloudv1.GetUserOptions)
				// Invoke operation with invalid model (negative test)
				result, response, operationErr = mqcloudService.GetUser(getUserOptionsModelNew)
				Expect(operationErr).ToNot(BeNil())
				Expect(response).To(BeNil())
				Expect(result).To(BeNil())
			})
			AfterEach(func() {
				testServer.Close()
			})
		})
		Context(`Using mock server endpoint with missing response body`, func() {
			BeforeEach(func() {
				testServer = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
					defer GinkgoRecover()

					// Set success status code with no respoonse body
					res.WriteHeader(200)
				}))
			})
			It(`Invoke GetUser successfully`, func() {
				mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1(&mqcloudv1.MqcloudV1Options{
					URL:            testServer.URL,
					Authenticator:  &core.NoAuthAuthenticator{},
					AcceptLanguage: core.StringPtr(acceptLanguage),
				})
				Expect(serviceErr).To(BeNil())
				Expect(mqcloudService).ToNot(BeNil())

				// Construct an instance of the GetUserOptions model
				getUserOptionsModel := new(mqcloudv1.GetUserOptions)
				getUserOptionsModel.ServiceInstanceGuid = core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				getUserOptionsModel.UserID = core.StringPtr("31a413dd84346effc8895b6ba4641641")
				getUserOptionsModel.Headers = map[string]string{"x-custom-header": "x-custom-value"}

				// Invoke operation
				result, response, operationErr := mqcloudService.GetUser(getUserOptionsModel)
				Expect(operationErr).To(BeNil())
				Expect(response).ToNot(BeNil())

				// Verify a nil result
				Expect(result).To(BeNil())
			})
			AfterEach(func() {
				testServer.Close()
			})
		})
	})
	Describe(`DeleteUser(deleteUserOptions *DeleteUserOptions)`, func() {
		acceptLanguage := "en-US,en;q=0.5"
		deleteUserPath := "/v1/a2b4d4bc-dadb-4637-bcec-9b7d1e723af8/users/31a413dd84346effc8895b6ba4641641"
		Context(`Using mock server endpoint`, func() {
			BeforeEach(func() {
				testServer = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
					defer GinkgoRecover()

					// Verify the contents of the request
					Expect(req.URL.EscapedPath()).To(Equal(deleteUserPath))
					Expect(req.Method).To(Equal("DELETE"))

					Expect(req.Header["Accept-Language"]).ToNot(BeNil())
					Expect(req.Header["Accept-Language"][0]).To(Equal(fmt.Sprintf("%v", "en-US,en;q=0.5")))
					res.WriteHeader(204)
				}))
			})
			It(`Invoke DeleteUser successfully`, func() {
				mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1(&mqcloudv1.MqcloudV1Options{
					URL:            testServer.URL,
					Authenticator:  &core.NoAuthAuthenticator{},
					AcceptLanguage: core.StringPtr(acceptLanguage),
				})
				Expect(serviceErr).To(BeNil())
				Expect(mqcloudService).ToNot(BeNil())

				// Invoke operation with nil options model (negative test)
				response, operationErr := mqcloudService.DeleteUser(nil)
				Expect(operationErr).NotTo(BeNil())
				Expect(response).To(BeNil())

				// Construct an instance of the DeleteUserOptions model
				deleteUserOptionsModel := new(mqcloudv1.DeleteUserOptions)
				deleteUserOptionsModel.ServiceInstanceGuid = core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				deleteUserOptionsModel.UserID = core.StringPtr("31a413dd84346effc8895b6ba4641641")
				deleteUserOptionsModel.Headers = map[string]string{"x-custom-header": "x-custom-value"}

				// Invoke operation with valid options model (positive test)
				response, operationErr = mqcloudService.DeleteUser(deleteUserOptionsModel)
				Expect(operationErr).To(BeNil())
				Expect(response).ToNot(BeNil())
			})
			It(`Invoke DeleteUser with error: Operation validation and request error`, func() {
				mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1(&mqcloudv1.MqcloudV1Options{
					URL:            testServer.URL,
					Authenticator:  &core.NoAuthAuthenticator{},
					AcceptLanguage: core.StringPtr(acceptLanguage),
				})
				Expect(serviceErr).To(BeNil())
				Expect(mqcloudService).ToNot(BeNil())

				// Construct an instance of the DeleteUserOptions model
				deleteUserOptionsModel := new(mqcloudv1.DeleteUserOptions)
				deleteUserOptionsModel.ServiceInstanceGuid = core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				deleteUserOptionsModel.UserID = core.StringPtr("31a413dd84346effc8895b6ba4641641")
				deleteUserOptionsModel.Headers = map[string]string{"x-custom-header": "x-custom-value"}
				// Invoke operation with empty URL (negative test)
				err := mqcloudService.SetServiceURL("")
				Expect(err).To(BeNil())
				response, operationErr := mqcloudService.DeleteUser(deleteUserOptionsModel)
				Expect(operationErr).ToNot(BeNil())
				Expect(operationErr.Error()).To(ContainSubstring(core.ERRORMSG_SERVICE_URL_MISSING))
				Expect(response).To(BeNil())
				// Construct a second instance of the DeleteUserOptions model with no property values
				deleteUserOptionsModelNew := new(mqcloudv1.DeleteUserOptions)
				// Invoke operation with invalid model (negative test)
				response, operationErr = mqcloudService.DeleteUser(deleteUserOptionsModelNew)
				Expect(operationErr).ToNot(BeNil())
				Expect(response).To(BeNil())
			})
			AfterEach(func() {
				testServer.Close()
			})
		})
	})
	Describe(`SetUserName(setUserNameOptions *SetUserNameOptions) - Operation response error`, func() {
		acceptLanguage := "en-US,en;q=0.5"
		setUserNamePath := "/v1/a2b4d4bc-dadb-4637-bcec-9b7d1e723af8/users/31a413dd84346effc8895b6ba4641641"
		Context(`Using mock server endpoint with invalid JSON response`, func() {
			BeforeEach(func() {
				testServer = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
					defer GinkgoRecover()

					// Verify the contents of the request
					Expect(req.URL.EscapedPath()).To(Equal(setUserNamePath))
					Expect(req.Method).To(Equal("PATCH"))
					Expect(req.Header["Accept-Language"]).ToNot(BeNil())
					Expect(req.Header["Accept-Language"][0]).To(Equal(fmt.Sprintf("%v", "en-US,en;q=0.5")))
					res.Header().Set("Content-type", "application/json")
					res.WriteHeader(200)
					fmt.Fprint(res, `} this is not valid json {`)
				}))
			})
			It(`Invoke SetUserName with error: Operation response processing error`, func() {
				mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1(&mqcloudv1.MqcloudV1Options{
					URL:            testServer.URL,
					Authenticator:  &core.NoAuthAuthenticator{},
					AcceptLanguage: core.StringPtr(acceptLanguage),
				})
				Expect(serviceErr).To(BeNil())
				Expect(mqcloudService).ToNot(BeNil())

				// Construct an instance of the SetUserNameOptions model
				setUserNameOptionsModel := new(mqcloudv1.SetUserNameOptions)
				setUserNameOptionsModel.ServiceInstanceGuid = core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				setUserNameOptionsModel.UserID = core.StringPtr("31a413dd84346effc8895b6ba4641641")
				setUserNameOptionsModel.Name = core.StringPtr("testString")
				setUserNameOptionsModel.Headers = map[string]string{"x-custom-header": "x-custom-value"}
				// Expect response parsing to fail since we are receiving a text/plain response
				result, response, operationErr := mqcloudService.SetUserName(setUserNameOptionsModel)
				Expect(operationErr).ToNot(BeNil())
				Expect(response).ToNot(BeNil())
				Expect(result).To(BeNil())

				// Enable retries and test again
				mqcloudService.EnableRetries(0, 0)
				result, response, operationErr = mqcloudService.SetUserName(setUserNameOptionsModel)
				Expect(operationErr).ToNot(BeNil())
				Expect(response).ToNot(BeNil())
				Expect(result).To(BeNil())
			})
			AfterEach(func() {
				testServer.Close()
			})
		})
	})
	Describe(`SetUserName(setUserNameOptions *SetUserNameOptions)`, func() {
		acceptLanguage := "en-US,en;q=0.5"
		setUserNamePath := "/v1/a2b4d4bc-dadb-4637-bcec-9b7d1e723af8/users/31a413dd84346effc8895b6ba4641641"
		Context(`Using mock server endpoint with timeout`, func() {
			BeforeEach(func() {
				testServer = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
					defer GinkgoRecover()

					// Verify the contents of the request
					Expect(req.URL.EscapedPath()).To(Equal(setUserNamePath))
					Expect(req.Method).To(Equal("PATCH"))

					// For gzip-disabled operation, verify Content-Encoding is not set.
					Expect(req.Header.Get("Content-Encoding")).To(BeEmpty())

					// If there is a body, then make sure we can read it
					bodyBuf := new(bytes.Buffer)
					if req.Header.Get("Content-Encoding") == "gzip" {
						body, err := core.NewGzipDecompressionReader(req.Body)
						Expect(err).To(BeNil())
						_, err = bodyBuf.ReadFrom(body)
						Expect(err).To(BeNil())
					} else {
						_, err := bodyBuf.ReadFrom(req.Body)
						Expect(err).To(BeNil())
					}
					fmt.Fprintf(GinkgoWriter, "  Request body: %s", bodyBuf.String())

					Expect(req.Header["Accept-Language"]).ToNot(BeNil())
					Expect(req.Header["Accept-Language"][0]).To(Equal(fmt.Sprintf("%v", "en-US,en;q=0.5")))
					// Sleep a short time to support a timeout test
					time.Sleep(100 * time.Millisecond)

					// Set mock response
					res.Header().Set("Content-type", "application/json")
					res.WriteHeader(200)
					fmt.Fprintf(res, "%s", `{"id": "ID", "name": "Name", "email": "user@host.org", "iam_service_id": "IamServiceID", "roles": ["Roles"], "iam_managed": true, "href": "Href"}`)
				}))
			})
			It(`Invoke SetUserName successfully with retries`, func() {
				mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1(&mqcloudv1.MqcloudV1Options{
					URL:            testServer.URL,
					Authenticator:  &core.NoAuthAuthenticator{},
					AcceptLanguage: core.StringPtr(acceptLanguage),
				})
				Expect(serviceErr).To(BeNil())
				Expect(mqcloudService).ToNot(BeNil())
				mqcloudService.EnableRetries(0, 0)

				// Construct an instance of the SetUserNameOptions model
				setUserNameOptionsModel := new(mqcloudv1.SetUserNameOptions)
				setUserNameOptionsModel.ServiceInstanceGuid = core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				setUserNameOptionsModel.UserID = core.StringPtr("31a413dd84346effc8895b6ba4641641")
				setUserNameOptionsModel.Name = core.StringPtr("testString")
				setUserNameOptionsModel.Headers = map[string]string{"x-custom-header": "x-custom-value"}

				// Invoke operation with a Context to test a timeout error
				ctx, cancelFunc := context.WithTimeout(context.Background(), 80*time.Millisecond)
				defer cancelFunc()
				_, _, operationErr := mqcloudService.SetUserNameWithContext(ctx, setUserNameOptionsModel)
				Expect(operationErr).ToNot(BeNil())
				Expect(operationErr.Error()).To(ContainSubstring("deadline exceeded"))

				// Disable retries and test again
				mqcloudService.DisableRetries()
				result, response, operationErr := mqcloudService.SetUserName(setUserNameOptionsModel)
				Expect(operationErr).To(BeNil())
				Expect(response).ToNot(BeNil())
				Expect(result).ToNot(BeNil())

				// Re-test the timeout error with retries disabled
				ctx, cancelFunc2 := context.WithTimeout(context.Background(), 80*time.Millisecond)
				defer cancelFunc2()
				_, _, operationErr = mqcloudService.SetUserNameWithContext(ctx, setUserNameOptionsModel)
				Expect(operationErr).ToNot(BeNil())
				Expect(operationErr.Error()).To(ContainSubstring("deadline exceeded"))
			})
			AfterEach(func() {
				testServer.Close()
			})
		})
		Context(`Using mock server endpoint`, func() {
			BeforeEach(func() {
				testServer = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
					defer GinkgoRecover()

					// Verify the contents of the request
					Expect(req.URL.EscapedPath()).To(Equal(setUserNamePath))
					Expect(req.Method).To(Equal("PATCH"))

					// For gzip-disabled operation, verify Content-Encoding is not set.
					Expect(req.Header.Get("Content-Encoding")).To(BeEmpty())

					// If there is a body, then make sure we can read it
					bodyBuf := new(bytes.Buffer)
					if req.Header.Get("Content-Encoding") == "gzip" {
						body, err := core.NewGzipDecompressionReader(req.Body)
						Expect(err).To(BeNil())
						_, err = bodyBuf.ReadFrom(body)
						Expect(err).To(BeNil())
					} else {
						_, err := bodyBuf.ReadFrom(req.Body)
						Expect(err).To(BeNil())
					}
					fmt.Fprintf(GinkgoWriter, "  Request body: %s", bodyBuf.String())

					Expect(req.Header["Accept-Language"]).ToNot(BeNil())
					Expect(req.Header["Accept-Language"][0]).To(Equal(fmt.Sprintf("%v", "en-US,en;q=0.5")))
					// Set mock response
					res.Header().Set("Content-type", "application/json")
					res.WriteHeader(200)
					fmt.Fprintf(res, "%s", `{"id": "ID", "name": "Name", "email": "user@host.org", "iam_service_id": "IamServiceID", "roles": ["Roles"], "iam_managed": true, "href": "Href"}`)
				}))
			})
			It(`Invoke SetUserName successfully`, func() {
				mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1(&mqcloudv1.MqcloudV1Options{
					URL:            testServer.URL,
					Authenticator:  &core.NoAuthAuthenticator{},
					AcceptLanguage: core.StringPtr(acceptLanguage),
				})
				Expect(serviceErr).To(BeNil())
				Expect(mqcloudService).ToNot(BeNil())

				// Invoke operation with nil options model (negative test)
				result, response, operationErr := mqcloudService.SetUserName(nil)
				Expect(operationErr).NotTo(BeNil())
				Expect(response).To(BeNil())
				Expect(result).To(BeNil())

				// Construct an instance of the SetUserNameOptions model
				setUserNameOptionsModel := new(mqcloudv1.SetUserNameOptions)
				setUserNameOptionsModel.ServiceInstanceGuid = core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				setUserNameOptionsModel.UserID = core.StringPtr("31a413dd84346effc8895b6ba4641641")
				setUserNameOptionsModel.Name = core.StringPtr("testString")
				setUserNameOptionsModel.Headers = map[string]string{"x-custom-header": "x-custom-value"}

				// Invoke operation with valid options model (positive test)
				result, response, operationErr = mqcloudService.SetUserName(setUserNameOptionsModel)
				Expect(operationErr).To(BeNil())
				Expect(response).ToNot(BeNil())
				Expect(result).ToNot(BeNil())

			})
			It(`Invoke SetUserName with error: Operation validation and request error`, func() {
				mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1(&mqcloudv1.MqcloudV1Options{
					URL:            testServer.URL,
					Authenticator:  &core.NoAuthAuthenticator{},
					AcceptLanguage: core.StringPtr(acceptLanguage),
				})
				Expect(serviceErr).To(BeNil())
				Expect(mqcloudService).ToNot(BeNil())

				// Construct an instance of the SetUserNameOptions model
				setUserNameOptionsModel := new(mqcloudv1.SetUserNameOptions)
				setUserNameOptionsModel.ServiceInstanceGuid = core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				setUserNameOptionsModel.UserID = core.StringPtr("31a413dd84346effc8895b6ba4641641")
				setUserNameOptionsModel.Name = core.StringPtr("testString")
				setUserNameOptionsModel.Headers = map[string]string{"x-custom-header": "x-custom-value"}
				// Invoke operation with empty URL (negative test)
				err := mqcloudService.SetServiceURL("")
				Expect(err).To(BeNil())
				result, response, operationErr := mqcloudService.SetUserName(setUserNameOptionsModel)
				Expect(operationErr).ToNot(BeNil())
				Expect(operationErr.Error()).To(ContainSubstring(core.ERRORMSG_SERVICE_URL_MISSING))
				Expect(response).To(BeNil())
				Expect(result).To(BeNil())
				// Construct a second instance of the SetUserNameOptions model with no property values
				setUserNameOptionsModelNew := new(mqcloudv1.SetUserNameOptions)
				// Invoke operation with invalid model (negative test)
				result, response, operationErr = mqcloudService.SetUserName(setUserNameOptionsModelNew)
				Expect(operationErr).ToNot(BeNil())
				Expect(response).To(BeNil())
				Expect(result).To(BeNil())
			})
			AfterEach(func() {
				testServer.Close()
			})
		})
		Context(`Using mock server endpoint with missing response body`, func() {
			BeforeEach(func() {
				testServer = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
					defer GinkgoRecover()

					// Set success status code with no respoonse body
					res.WriteHeader(200)
				}))
			})
			It(`Invoke SetUserName successfully`, func() {
				mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1(&mqcloudv1.MqcloudV1Options{
					URL:            testServer.URL,
					Authenticator:  &core.NoAuthAuthenticator{},
					AcceptLanguage: core.StringPtr(acceptLanguage),
				})
				Expect(serviceErr).To(BeNil())
				Expect(mqcloudService).ToNot(BeNil())

				// Construct an instance of the SetUserNameOptions model
				setUserNameOptionsModel := new(mqcloudv1.SetUserNameOptions)
				setUserNameOptionsModel.ServiceInstanceGuid = core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				setUserNameOptionsModel.UserID = core.StringPtr("31a413dd84346effc8895b6ba4641641")
				setUserNameOptionsModel.Name = core.StringPtr("testString")
				setUserNameOptionsModel.Headers = map[string]string{"x-custom-header": "x-custom-value"}

				// Invoke operation
				result, response, operationErr := mqcloudService.SetUserName(setUserNameOptionsModel)
				Expect(operationErr).To(BeNil())
				Expect(response).ToNot(BeNil())

				// Verify a nil result
				Expect(result).To(BeNil())
			})
			AfterEach(func() {
				testServer.Close()
			})
		})
	})
	Describe(`ListApplications(listApplicationsOptions *ListApplicationsOptions) - Operation response error`, func() {
		acceptLanguage := "en-US,en;q=0.5"
		listApplicationsPath := "/v1/a2b4d4bc-dadb-4637-bcec-9b7d1e723af8/applications"
		Context(`Using mock server endpoint with invalid JSON response`, func() {
			BeforeEach(func() {
				testServer = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
					defer GinkgoRecover()

					// Verify the contents of the request
					Expect(req.URL.EscapedPath()).To(Equal(listApplicationsPath))
					Expect(req.Method).To(Equal("GET"))
					Expect(req.Header["Accept-Language"]).ToNot(BeNil())
					Expect(req.Header["Accept-Language"][0]).To(Equal(fmt.Sprintf("%v", "en-US,en;q=0.5")))
					Expect(req.URL.Query()["offset"]).To(Equal([]string{fmt.Sprint(int64(0))}))
					Expect(req.URL.Query()["limit"]).To(Equal([]string{fmt.Sprint(int64(10))}))
					res.Header().Set("Content-type", "application/json")
					res.WriteHeader(200)
					fmt.Fprint(res, `} this is not valid json {`)
				}))
			})
			It(`Invoke ListApplications with error: Operation response processing error`, func() {
				mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1(&mqcloudv1.MqcloudV1Options{
					URL:            testServer.URL,
					Authenticator:  &core.NoAuthAuthenticator{},
					AcceptLanguage: core.StringPtr(acceptLanguage),
				})
				Expect(serviceErr).To(BeNil())
				Expect(mqcloudService).ToNot(BeNil())

				// Construct an instance of the ListApplicationsOptions model
				listApplicationsOptionsModel := new(mqcloudv1.ListApplicationsOptions)
				listApplicationsOptionsModel.ServiceInstanceGuid = core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				listApplicationsOptionsModel.Offset = core.Int64Ptr(int64(0))
				listApplicationsOptionsModel.Limit = core.Int64Ptr(int64(10))
				listApplicationsOptionsModel.Headers = map[string]string{"x-custom-header": "x-custom-value"}
				// Expect response parsing to fail since we are receiving a text/plain response
				result, response, operationErr := mqcloudService.ListApplications(listApplicationsOptionsModel)
				Expect(operationErr).ToNot(BeNil())
				Expect(response).ToNot(BeNil())
				Expect(result).To(BeNil())

				// Enable retries and test again
				mqcloudService.EnableRetries(0, 0)
				result, response, operationErr = mqcloudService.ListApplications(listApplicationsOptionsModel)
				Expect(operationErr).ToNot(BeNil())
				Expect(response).ToNot(BeNil())
				Expect(result).To(BeNil())
			})
			AfterEach(func() {
				testServer.Close()
			})
		})
	})
	Describe(`ListApplications(listApplicationsOptions *ListApplicationsOptions)`, func() {
		acceptLanguage := "en-US,en;q=0.5"
		listApplicationsPath := "/v1/a2b4d4bc-dadb-4637-bcec-9b7d1e723af8/applications"
		Context(`Using mock server endpoint with timeout`, func() {
			BeforeEach(func() {
				testServer = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
					defer GinkgoRecover()

					// Verify the contents of the request
					Expect(req.URL.EscapedPath()).To(Equal(listApplicationsPath))
					Expect(req.Method).To(Equal("GET"))

					Expect(req.Header["Accept-Language"]).ToNot(BeNil())
					Expect(req.Header["Accept-Language"][0]).To(Equal(fmt.Sprintf("%v", "en-US,en;q=0.5")))
					Expect(req.URL.Query()["offset"]).To(Equal([]string{fmt.Sprint(int64(0))}))
					Expect(req.URL.Query()["limit"]).To(Equal([]string{fmt.Sprint(int64(10))}))
					// Sleep a short time to support a timeout test
					time.Sleep(100 * time.Millisecond)

					// Set mock response
					res.Header().Set("Content-type", "application/json")
					res.WriteHeader(200)
					fmt.Fprintf(res, "%s", `{"offset": 6, "limit": 50, "first": {"href": "Href"}, "next": {"href": "Href"}, "previous": {"href": "Href"}, "applications": [{"id": "ID", "name": "Name", "iam_service_id": "IamServiceID", "create_api_key_uri": "CreateApiKeyURI", "href": "Href"}]}`)
				}))
			})
			It(`Invoke ListApplications successfully with retries`, func() {
				mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1(&mqcloudv1.MqcloudV1Options{
					URL:            testServer.URL,
					Authenticator:  &core.NoAuthAuthenticator{},
					AcceptLanguage: core.StringPtr(acceptLanguage),
				})
				Expect(serviceErr).To(BeNil())
				Expect(mqcloudService).ToNot(BeNil())
				mqcloudService.EnableRetries(0, 0)

				// Construct an instance of the ListApplicationsOptions model
				listApplicationsOptionsModel := new(mqcloudv1.ListApplicationsOptions)
				listApplicationsOptionsModel.ServiceInstanceGuid = core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				listApplicationsOptionsModel.Offset = core.Int64Ptr(int64(0))
				listApplicationsOptionsModel.Limit = core.Int64Ptr(int64(10))
				listApplicationsOptionsModel.Headers = map[string]string{"x-custom-header": "x-custom-value"}

				// Invoke operation with a Context to test a timeout error
				ctx, cancelFunc := context.WithTimeout(context.Background(), 80*time.Millisecond)
				defer cancelFunc()
				_, _, operationErr := mqcloudService.ListApplicationsWithContext(ctx, listApplicationsOptionsModel)
				Expect(operationErr).ToNot(BeNil())
				Expect(operationErr.Error()).To(ContainSubstring("deadline exceeded"))

				// Disable retries and test again
				mqcloudService.DisableRetries()
				result, response, operationErr := mqcloudService.ListApplications(listApplicationsOptionsModel)
				Expect(operationErr).To(BeNil())
				Expect(response).ToNot(BeNil())
				Expect(result).ToNot(BeNil())

				// Re-test the timeout error with retries disabled
				ctx, cancelFunc2 := context.WithTimeout(context.Background(), 80*time.Millisecond)
				defer cancelFunc2()
				_, _, operationErr = mqcloudService.ListApplicationsWithContext(ctx, listApplicationsOptionsModel)
				Expect(operationErr).ToNot(BeNil())
				Expect(operationErr.Error()).To(ContainSubstring("deadline exceeded"))
			})
			AfterEach(func() {
				testServer.Close()
			})
		})
		Context(`Using mock server endpoint`, func() {
			BeforeEach(func() {
				testServer = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
					defer GinkgoRecover()

					// Verify the contents of the request
					Expect(req.URL.EscapedPath()).To(Equal(listApplicationsPath))
					Expect(req.Method).To(Equal("GET"))

					Expect(req.Header["Accept-Language"]).ToNot(BeNil())
					Expect(req.Header["Accept-Language"][0]).To(Equal(fmt.Sprintf("%v", "en-US,en;q=0.5")))
					Expect(req.URL.Query()["offset"]).To(Equal([]string{fmt.Sprint(int64(0))}))
					Expect(req.URL.Query()["limit"]).To(Equal([]string{fmt.Sprint(int64(10))}))
					// Set mock response
					res.Header().Set("Content-type", "application/json")
					res.WriteHeader(200)
					fmt.Fprintf(res, "%s", `{"offset": 6, "limit": 50, "first": {"href": "Href"}, "next": {"href": "Href"}, "previous": {"href": "Href"}, "applications": [{"id": "ID", "name": "Name", "iam_service_id": "IamServiceID", "create_api_key_uri": "CreateApiKeyURI", "href": "Href"}]}`)
				}))
			})
			It(`Invoke ListApplications successfully`, func() {
				mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1(&mqcloudv1.MqcloudV1Options{
					URL:            testServer.URL,
					Authenticator:  &core.NoAuthAuthenticator{},
					AcceptLanguage: core.StringPtr(acceptLanguage),
				})
				Expect(serviceErr).To(BeNil())
				Expect(mqcloudService).ToNot(BeNil())

				// Invoke operation with nil options model (negative test)
				result, response, operationErr := mqcloudService.ListApplications(nil)
				Expect(operationErr).NotTo(BeNil())
				Expect(response).To(BeNil())
				Expect(result).To(BeNil())

				// Construct an instance of the ListApplicationsOptions model
				listApplicationsOptionsModel := new(mqcloudv1.ListApplicationsOptions)
				listApplicationsOptionsModel.ServiceInstanceGuid = core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				listApplicationsOptionsModel.Offset = core.Int64Ptr(int64(0))
				listApplicationsOptionsModel.Limit = core.Int64Ptr(int64(10))
				listApplicationsOptionsModel.Headers = map[string]string{"x-custom-header": "x-custom-value"}

				// Invoke operation with valid options model (positive test)
				result, response, operationErr = mqcloudService.ListApplications(listApplicationsOptionsModel)
				Expect(operationErr).To(BeNil())
				Expect(response).ToNot(BeNil())
				Expect(result).ToNot(BeNil())

			})
			It(`Invoke ListApplications with error: Operation validation and request error`, func() {
				mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1(&mqcloudv1.MqcloudV1Options{
					URL:            testServer.URL,
					Authenticator:  &core.NoAuthAuthenticator{},
					AcceptLanguage: core.StringPtr(acceptLanguage),
				})
				Expect(serviceErr).To(BeNil())
				Expect(mqcloudService).ToNot(BeNil())

				// Construct an instance of the ListApplicationsOptions model
				listApplicationsOptionsModel := new(mqcloudv1.ListApplicationsOptions)
				listApplicationsOptionsModel.ServiceInstanceGuid = core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				listApplicationsOptionsModel.Offset = core.Int64Ptr(int64(0))
				listApplicationsOptionsModel.Limit = core.Int64Ptr(int64(10))
				listApplicationsOptionsModel.Headers = map[string]string{"x-custom-header": "x-custom-value"}
				// Invoke operation with empty URL (negative test)
				err := mqcloudService.SetServiceURL("")
				Expect(err).To(BeNil())
				result, response, operationErr := mqcloudService.ListApplications(listApplicationsOptionsModel)
				Expect(operationErr).ToNot(BeNil())
				Expect(operationErr.Error()).To(ContainSubstring(core.ERRORMSG_SERVICE_URL_MISSING))
				Expect(response).To(BeNil())
				Expect(result).To(BeNil())
				// Construct a second instance of the ListApplicationsOptions model with no property values
				listApplicationsOptionsModelNew := new(mqcloudv1.ListApplicationsOptions)
				// Invoke operation with invalid model (negative test)
				result, response, operationErr = mqcloudService.ListApplications(listApplicationsOptionsModelNew)
				Expect(operationErr).ToNot(BeNil())
				Expect(response).To(BeNil())
				Expect(result).To(BeNil())
			})
			AfterEach(func() {
				testServer.Close()
			})
		})
		Context(`Using mock server endpoint with missing response body`, func() {
			BeforeEach(func() {
				testServer = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
					defer GinkgoRecover()

					// Set success status code with no respoonse body
					res.WriteHeader(200)
				}))
			})
			It(`Invoke ListApplications successfully`, func() {
				mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1(&mqcloudv1.MqcloudV1Options{
					URL:            testServer.URL,
					Authenticator:  &core.NoAuthAuthenticator{},
					AcceptLanguage: core.StringPtr(acceptLanguage),
				})
				Expect(serviceErr).To(BeNil())
				Expect(mqcloudService).ToNot(BeNil())

				// Construct an instance of the ListApplicationsOptions model
				listApplicationsOptionsModel := new(mqcloudv1.ListApplicationsOptions)
				listApplicationsOptionsModel.ServiceInstanceGuid = core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				listApplicationsOptionsModel.Offset = core.Int64Ptr(int64(0))
				listApplicationsOptionsModel.Limit = core.Int64Ptr(int64(10))
				listApplicationsOptionsModel.Headers = map[string]string{"x-custom-header": "x-custom-value"}

				// Invoke operation
				result, response, operationErr := mqcloudService.ListApplications(listApplicationsOptionsModel)
				Expect(operationErr).To(BeNil())
				Expect(response).ToNot(BeNil())

				// Verify a nil result
				Expect(result).To(BeNil())
			})
			AfterEach(func() {
				testServer.Close()
			})
		})
		Context(`Test pagination helper method on response`, func() {
			It(`Invoke GetNextOffset successfully`, func() {
				responseObject := new(mqcloudv1.ApplicationDetailsCollection)
				nextObject := new(mqcloudv1.Next)
				nextObject.Href = core.StringPtr("ibm.com?offset=135")
				responseObject.Next = nextObject

				value, err := responseObject.GetNextOffset()
				Expect(err).To(BeNil())
				Expect(value).To(Equal(core.Int64Ptr(int64(135))))
			})
			It(`Invoke GetNextOffset without a "Next" property in the response`, func() {
				responseObject := new(mqcloudv1.ApplicationDetailsCollection)

				value, err := responseObject.GetNextOffset()
				Expect(err).To(BeNil())
				Expect(value).To(BeNil())
			})
			It(`Invoke GetNextOffset without any query params in the "Next" URL`, func() {
				responseObject := new(mqcloudv1.ApplicationDetailsCollection)
				nextObject := new(mqcloudv1.Next)
				nextObject.Href = core.StringPtr("ibm.com")
				responseObject.Next = nextObject

				value, err := responseObject.GetNextOffset()
				Expect(err).To(BeNil())
				Expect(value).To(BeNil())
			})
			It(`Invoke GetNextOffset with a non-integer query param in the "Next" URL`, func() {
				responseObject := new(mqcloudv1.ApplicationDetailsCollection)
				nextObject := new(mqcloudv1.Next)
				nextObject.Href = core.StringPtr("ibm.com?offset=tiger")
				responseObject.Next = nextObject

				value, err := responseObject.GetNextOffset()
				Expect(err).NotTo(BeNil())
				Expect(value).To(BeNil())
			})
		})
		Context(`Using mock server endpoint - paginated response`, func() {
			BeforeEach(func() {
				var requestNumber int = 0
				testServer = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
					defer GinkgoRecover()

					// Verify the contents of the request
					Expect(req.URL.EscapedPath()).To(Equal(listApplicationsPath))
					Expect(req.Method).To(Equal("GET"))

					// Set mock response
					res.Header().Set("Content-type", "application/json")
					res.WriteHeader(200)
					requestNumber++
					if requestNumber == 1 {
						fmt.Fprintf(res, "%s", `{"next":{"href":"https://myhost.com/somePath?offset=1"},"total_count":2,"limit":1,"applications":[{"id":"ID","name":"Name","iam_service_id":"IamServiceID","create_api_key_uri":"CreateApiKeyURI","href":"Href"}]}`)
					} else if requestNumber == 2 {
						fmt.Fprintf(res, "%s", `{"total_count":2,"limit":1,"applications":[{"id":"ID","name":"Name","iam_service_id":"IamServiceID","create_api_key_uri":"CreateApiKeyURI","href":"Href"}]}`)
					} else {
						res.WriteHeader(400)
					}
				}))
			})
			It(`Use ApplicationsPager.GetNext successfully`, func() {
				mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1(&mqcloudv1.MqcloudV1Options{
					URL:            testServer.URL,
					Authenticator:  &core.NoAuthAuthenticator{},
					AcceptLanguage: core.StringPtr(acceptLanguage),
				})
				Expect(serviceErr).To(BeNil())
				Expect(mqcloudService).ToNot(BeNil())

				listApplicationsOptionsModel := &mqcloudv1.ListApplicationsOptions{
					ServiceInstanceGuid: core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8"),
					Limit:               core.Int64Ptr(int64(10)),
				}

				pager, err := mqcloudService.NewApplicationsPager(listApplicationsOptionsModel)
				Expect(err).To(BeNil())
				Expect(pager).ToNot(BeNil())

				var allResults []mqcloudv1.ApplicationDetails
				for pager.HasNext() {
					nextPage, err := pager.GetNext()
					Expect(err).To(BeNil())
					Expect(nextPage).ToNot(BeNil())
					allResults = append(allResults, nextPage...)
				}
				Expect(len(allResults)).To(Equal(2))
			})
			It(`Use ApplicationsPager.GetAll successfully`, func() {
				mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1(&mqcloudv1.MqcloudV1Options{
					URL:            testServer.URL,
					Authenticator:  &core.NoAuthAuthenticator{},
					AcceptLanguage: core.StringPtr(acceptLanguage),
				})
				Expect(serviceErr).To(BeNil())
				Expect(mqcloudService).ToNot(BeNil())

				listApplicationsOptionsModel := &mqcloudv1.ListApplicationsOptions{
					ServiceInstanceGuid: core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8"),
					Limit:               core.Int64Ptr(int64(10)),
				}

				pager, err := mqcloudService.NewApplicationsPager(listApplicationsOptionsModel)
				Expect(err).To(BeNil())
				Expect(pager).ToNot(BeNil())

				allResults, err := pager.GetAll()
				Expect(err).To(BeNil())
				Expect(allResults).ToNot(BeNil())
				Expect(len(allResults)).To(Equal(2))
			})
		})
	})
	Describe(`CreateApplication(createApplicationOptions *CreateApplicationOptions) - Operation response error`, func() {
		acceptLanguage := "en-US,en;q=0.5"
		createApplicationPath := "/v1/a2b4d4bc-dadb-4637-bcec-9b7d1e723af8/applications"
		Context(`Using mock server endpoint with invalid JSON response`, func() {
			BeforeEach(func() {
				testServer = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
					defer GinkgoRecover()

					// Verify the contents of the request
					Expect(req.URL.EscapedPath()).To(Equal(createApplicationPath))
					Expect(req.Method).To(Equal("POST"))
					Expect(req.Header["Accept-Language"]).ToNot(BeNil())
					Expect(req.Header["Accept-Language"][0]).To(Equal(fmt.Sprintf("%v", "en-US,en;q=0.5")))
					res.Header().Set("Content-type", "application/json")
					res.WriteHeader(201)
					fmt.Fprint(res, `} this is not valid json {`)
				}))
			})
			It(`Invoke CreateApplication with error: Operation response processing error`, func() {
				mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1(&mqcloudv1.MqcloudV1Options{
					URL:            testServer.URL,
					Authenticator:  &core.NoAuthAuthenticator{},
					AcceptLanguage: core.StringPtr(acceptLanguage),
				})
				Expect(serviceErr).To(BeNil())
				Expect(mqcloudService).ToNot(BeNil())

				// Construct an instance of the CreateApplicationOptions model
				createApplicationOptionsModel := new(mqcloudv1.CreateApplicationOptions)
				createApplicationOptionsModel.ServiceInstanceGuid = core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				createApplicationOptionsModel.Name = core.StringPtr("test-app")
				createApplicationOptionsModel.Headers = map[string]string{"x-custom-header": "x-custom-value"}
				// Expect response parsing to fail since we are receiving a text/plain response
				result, response, operationErr := mqcloudService.CreateApplication(createApplicationOptionsModel)
				Expect(operationErr).ToNot(BeNil())
				Expect(response).ToNot(BeNil())
				Expect(result).To(BeNil())

				// Enable retries and test again
				mqcloudService.EnableRetries(0, 0)
				result, response, operationErr = mqcloudService.CreateApplication(createApplicationOptionsModel)
				Expect(operationErr).ToNot(BeNil())
				Expect(response).ToNot(BeNil())
				Expect(result).To(BeNil())
			})
			AfterEach(func() {
				testServer.Close()
			})
		})
	})
	Describe(`CreateApplication(createApplicationOptions *CreateApplicationOptions)`, func() {
		acceptLanguage := "en-US,en;q=0.5"
		createApplicationPath := "/v1/a2b4d4bc-dadb-4637-bcec-9b7d1e723af8/applications"
		Context(`Using mock server endpoint with timeout`, func() {
			BeforeEach(func() {
				testServer = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
					defer GinkgoRecover()

					// Verify the contents of the request
					Expect(req.URL.EscapedPath()).To(Equal(createApplicationPath))
					Expect(req.Method).To(Equal("POST"))

					// For gzip-disabled operation, verify Content-Encoding is not set.
					Expect(req.Header.Get("Content-Encoding")).To(BeEmpty())

					// If there is a body, then make sure we can read it
					bodyBuf := new(bytes.Buffer)
					if req.Header.Get("Content-Encoding") == "gzip" {
						body, err := core.NewGzipDecompressionReader(req.Body)
						Expect(err).To(BeNil())
						_, err = bodyBuf.ReadFrom(body)
						Expect(err).To(BeNil())
					} else {
						_, err := bodyBuf.ReadFrom(req.Body)
						Expect(err).To(BeNil())
					}
					fmt.Fprintf(GinkgoWriter, "  Request body: %s", bodyBuf.String())

					Expect(req.Header["Accept-Language"]).ToNot(BeNil())
					Expect(req.Header["Accept-Language"][0]).To(Equal(fmt.Sprintf("%v", "en-US,en;q=0.5")))
					// Sleep a short time to support a timeout test
					time.Sleep(100 * time.Millisecond)

					// Set mock response
					res.Header().Set("Content-type", "application/json")
					res.WriteHeader(201)
					fmt.Fprintf(res, "%s", `{"id": "ID", "name": "Name", "iam_service_id": "IamServiceID", "create_api_key_uri": "CreateApiKeyURI", "href": "Href", "api_key_name": "ApiKeyName", "api_key_id": "ApiKeyID", "api_key": "ApiKey"}`)
				}))
			})
			It(`Invoke CreateApplication successfully with retries`, func() {
				mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1(&mqcloudv1.MqcloudV1Options{
					URL:            testServer.URL,
					Authenticator:  &core.NoAuthAuthenticator{},
					AcceptLanguage: core.StringPtr(acceptLanguage),
				})
				Expect(serviceErr).To(BeNil())
				Expect(mqcloudService).ToNot(BeNil())
				mqcloudService.EnableRetries(0, 0)

				// Construct an instance of the CreateApplicationOptions model
				createApplicationOptionsModel := new(mqcloudv1.CreateApplicationOptions)
				createApplicationOptionsModel.ServiceInstanceGuid = core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				createApplicationOptionsModel.Name = core.StringPtr("test-app")
				createApplicationOptionsModel.Headers = map[string]string{"x-custom-header": "x-custom-value"}

				// Invoke operation with a Context to test a timeout error
				ctx, cancelFunc := context.WithTimeout(context.Background(), 80*time.Millisecond)
				defer cancelFunc()
				_, _, operationErr := mqcloudService.CreateApplicationWithContext(ctx, createApplicationOptionsModel)
				Expect(operationErr).ToNot(BeNil())
				Expect(operationErr.Error()).To(ContainSubstring("deadline exceeded"))

				// Disable retries and test again
				mqcloudService.DisableRetries()
				result, response, operationErr := mqcloudService.CreateApplication(createApplicationOptionsModel)
				Expect(operationErr).To(BeNil())
				Expect(response).ToNot(BeNil())
				Expect(result).ToNot(BeNil())

				// Re-test the timeout error with retries disabled
				ctx, cancelFunc2 := context.WithTimeout(context.Background(), 80*time.Millisecond)
				defer cancelFunc2()
				_, _, operationErr = mqcloudService.CreateApplicationWithContext(ctx, createApplicationOptionsModel)
				Expect(operationErr).ToNot(BeNil())
				Expect(operationErr.Error()).To(ContainSubstring("deadline exceeded"))
			})
			AfterEach(func() {
				testServer.Close()
			})
		})
		Context(`Using mock server endpoint`, func() {
			BeforeEach(func() {
				testServer = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
					defer GinkgoRecover()

					// Verify the contents of the request
					Expect(req.URL.EscapedPath()).To(Equal(createApplicationPath))
					Expect(req.Method).To(Equal("POST"))

					// For gzip-disabled operation, verify Content-Encoding is not set.
					Expect(req.Header.Get("Content-Encoding")).To(BeEmpty())

					// If there is a body, then make sure we can read it
					bodyBuf := new(bytes.Buffer)
					if req.Header.Get("Content-Encoding") == "gzip" {
						body, err := core.NewGzipDecompressionReader(req.Body)
						Expect(err).To(BeNil())
						_, err = bodyBuf.ReadFrom(body)
						Expect(err).To(BeNil())
					} else {
						_, err := bodyBuf.ReadFrom(req.Body)
						Expect(err).To(BeNil())
					}
					fmt.Fprintf(GinkgoWriter, "  Request body: %s", bodyBuf.String())

					Expect(req.Header["Accept-Language"]).ToNot(BeNil())
					Expect(req.Header["Accept-Language"][0]).To(Equal(fmt.Sprintf("%v", "en-US,en;q=0.5")))
					// Set mock response
					res.Header().Set("Content-type", "application/json")
					res.WriteHeader(201)
					fmt.Fprintf(res, "%s", `{"id": "ID", "name": "Name", "iam_service_id": "IamServiceID", "create_api_key_uri": "CreateApiKeyURI", "href": "Href", "api_key_name": "ApiKeyName", "api_key_id": "ApiKeyID", "api_key": "ApiKey"}`)
				}))
			})
			It(`Invoke CreateApplication successfully`, func() {
				mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1(&mqcloudv1.MqcloudV1Options{
					URL:            testServer.URL,
					Authenticator:  &core.NoAuthAuthenticator{},
					AcceptLanguage: core.StringPtr(acceptLanguage),
				})
				Expect(serviceErr).To(BeNil())
				Expect(mqcloudService).ToNot(BeNil())

				// Invoke operation with nil options model (negative test)
				result, response, operationErr := mqcloudService.CreateApplication(nil)
				Expect(operationErr).NotTo(BeNil())
				Expect(response).To(BeNil())
				Expect(result).To(BeNil())

				// Construct an instance of the CreateApplicationOptions model
				createApplicationOptionsModel := new(mqcloudv1.CreateApplicationOptions)
				createApplicationOptionsModel.ServiceInstanceGuid = core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				createApplicationOptionsModel.Name = core.StringPtr("test-app")
				createApplicationOptionsModel.Headers = map[string]string{"x-custom-header": "x-custom-value"}

				// Invoke operation with valid options model (positive test)
				result, response, operationErr = mqcloudService.CreateApplication(createApplicationOptionsModel)
				Expect(operationErr).To(BeNil())
				Expect(response).ToNot(BeNil())
				Expect(result).ToNot(BeNil())

			})
			It(`Invoke CreateApplication with error: Operation validation and request error`, func() {
				mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1(&mqcloudv1.MqcloudV1Options{
					URL:            testServer.URL,
					Authenticator:  &core.NoAuthAuthenticator{},
					AcceptLanguage: core.StringPtr(acceptLanguage),
				})
				Expect(serviceErr).To(BeNil())
				Expect(mqcloudService).ToNot(BeNil())

				// Construct an instance of the CreateApplicationOptions model
				createApplicationOptionsModel := new(mqcloudv1.CreateApplicationOptions)
				createApplicationOptionsModel.ServiceInstanceGuid = core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				createApplicationOptionsModel.Name = core.StringPtr("test-app")
				createApplicationOptionsModel.Headers = map[string]string{"x-custom-header": "x-custom-value"}
				// Invoke operation with empty URL (negative test)
				err := mqcloudService.SetServiceURL("")
				Expect(err).To(BeNil())
				result, response, operationErr := mqcloudService.CreateApplication(createApplicationOptionsModel)
				Expect(operationErr).ToNot(BeNil())
				Expect(operationErr.Error()).To(ContainSubstring(core.ERRORMSG_SERVICE_URL_MISSING))
				Expect(response).To(BeNil())
				Expect(result).To(BeNil())
				// Construct a second instance of the CreateApplicationOptions model with no property values
				createApplicationOptionsModelNew := new(mqcloudv1.CreateApplicationOptions)
				// Invoke operation with invalid model (negative test)
				result, response, operationErr = mqcloudService.CreateApplication(createApplicationOptionsModelNew)
				Expect(operationErr).ToNot(BeNil())
				Expect(response).To(BeNil())
				Expect(result).To(BeNil())
			})
			AfterEach(func() {
				testServer.Close()
			})
		})
		Context(`Using mock server endpoint with missing response body`, func() {
			BeforeEach(func() {
				testServer = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
					defer GinkgoRecover()

					// Set success status code with no respoonse body
					res.WriteHeader(201)
				}))
			})
			It(`Invoke CreateApplication successfully`, func() {
				mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1(&mqcloudv1.MqcloudV1Options{
					URL:            testServer.URL,
					Authenticator:  &core.NoAuthAuthenticator{},
					AcceptLanguage: core.StringPtr(acceptLanguage),
				})
				Expect(serviceErr).To(BeNil())
				Expect(mqcloudService).ToNot(BeNil())

				// Construct an instance of the CreateApplicationOptions model
				createApplicationOptionsModel := new(mqcloudv1.CreateApplicationOptions)
				createApplicationOptionsModel.ServiceInstanceGuid = core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				createApplicationOptionsModel.Name = core.StringPtr("test-app")
				createApplicationOptionsModel.Headers = map[string]string{"x-custom-header": "x-custom-value"}

				// Invoke operation
				result, response, operationErr := mqcloudService.CreateApplication(createApplicationOptionsModel)
				Expect(operationErr).To(BeNil())
				Expect(response).ToNot(BeNil())

				// Verify a nil result
				Expect(result).To(BeNil())
			})
			AfterEach(func() {
				testServer.Close()
			})
		})
	})
	Describe(`GetApplication(getApplicationOptions *GetApplicationOptions) - Operation response error`, func() {
		acceptLanguage := "en-US,en;q=0.5"
		getApplicationPath := "/v1/a2b4d4bc-dadb-4637-bcec-9b7d1e723af8/applications/0123456789ABCDEF0123456789ABCDEF"
		Context(`Using mock server endpoint with invalid JSON response`, func() {
			BeforeEach(func() {
				testServer = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
					defer GinkgoRecover()

					// Verify the contents of the request
					Expect(req.URL.EscapedPath()).To(Equal(getApplicationPath))
					Expect(req.Method).To(Equal("GET"))
					Expect(req.Header["Accept-Language"]).ToNot(BeNil())
					Expect(req.Header["Accept-Language"][0]).To(Equal(fmt.Sprintf("%v", "en-US,en;q=0.5")))
					res.Header().Set("Content-type", "application/json")
					res.WriteHeader(200)
					fmt.Fprint(res, `} this is not valid json {`)
				}))
			})
			It(`Invoke GetApplication with error: Operation response processing error`, func() {
				mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1(&mqcloudv1.MqcloudV1Options{
					URL:            testServer.URL,
					Authenticator:  &core.NoAuthAuthenticator{},
					AcceptLanguage: core.StringPtr(acceptLanguage),
				})
				Expect(serviceErr).To(BeNil())
				Expect(mqcloudService).ToNot(BeNil())

				// Construct an instance of the GetApplicationOptions model
				getApplicationOptionsModel := new(mqcloudv1.GetApplicationOptions)
				getApplicationOptionsModel.ServiceInstanceGuid = core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				getApplicationOptionsModel.ApplicationID = core.StringPtr("0123456789ABCDEF0123456789ABCDEF")
				getApplicationOptionsModel.Headers = map[string]string{"x-custom-header": "x-custom-value"}
				// Expect response parsing to fail since we are receiving a text/plain response
				result, response, operationErr := mqcloudService.GetApplication(getApplicationOptionsModel)
				Expect(operationErr).ToNot(BeNil())
				Expect(response).ToNot(BeNil())
				Expect(result).To(BeNil())

				// Enable retries and test again
				mqcloudService.EnableRetries(0, 0)
				result, response, operationErr = mqcloudService.GetApplication(getApplicationOptionsModel)
				Expect(operationErr).ToNot(BeNil())
				Expect(response).ToNot(BeNil())
				Expect(result).To(BeNil())
			})
			AfterEach(func() {
				testServer.Close()
			})
		})
	})
	Describe(`GetApplication(getApplicationOptions *GetApplicationOptions)`, func() {
		acceptLanguage := "en-US,en;q=0.5"
		getApplicationPath := "/v1/a2b4d4bc-dadb-4637-bcec-9b7d1e723af8/applications/0123456789ABCDEF0123456789ABCDEF"
		Context(`Using mock server endpoint with timeout`, func() {
			BeforeEach(func() {
				testServer = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
					defer GinkgoRecover()

					// Verify the contents of the request
					Expect(req.URL.EscapedPath()).To(Equal(getApplicationPath))
					Expect(req.Method).To(Equal("GET"))

					Expect(req.Header["Accept-Language"]).ToNot(BeNil())
					Expect(req.Header["Accept-Language"][0]).To(Equal(fmt.Sprintf("%v", "en-US,en;q=0.5")))
					// Sleep a short time to support a timeout test
					time.Sleep(100 * time.Millisecond)

					// Set mock response
					res.Header().Set("Content-type", "application/json")
					res.WriteHeader(200)
					fmt.Fprintf(res, "%s", `{"id": "ID", "name": "Name", "iam_service_id": "IamServiceID", "create_api_key_uri": "CreateApiKeyURI", "href": "Href"}`)
				}))
			})
			It(`Invoke GetApplication successfully with retries`, func() {
				mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1(&mqcloudv1.MqcloudV1Options{
					URL:            testServer.URL,
					Authenticator:  &core.NoAuthAuthenticator{},
					AcceptLanguage: core.StringPtr(acceptLanguage),
				})
				Expect(serviceErr).To(BeNil())
				Expect(mqcloudService).ToNot(BeNil())
				mqcloudService.EnableRetries(0, 0)

				// Construct an instance of the GetApplicationOptions model
				getApplicationOptionsModel := new(mqcloudv1.GetApplicationOptions)
				getApplicationOptionsModel.ServiceInstanceGuid = core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				getApplicationOptionsModel.ApplicationID = core.StringPtr("0123456789ABCDEF0123456789ABCDEF")
				getApplicationOptionsModel.Headers = map[string]string{"x-custom-header": "x-custom-value"}

				// Invoke operation with a Context to test a timeout error
				ctx, cancelFunc := context.WithTimeout(context.Background(), 80*time.Millisecond)
				defer cancelFunc()
				_, _, operationErr := mqcloudService.GetApplicationWithContext(ctx, getApplicationOptionsModel)
				Expect(operationErr).ToNot(BeNil())
				Expect(operationErr.Error()).To(ContainSubstring("deadline exceeded"))

				// Disable retries and test again
				mqcloudService.DisableRetries()
				result, response, operationErr := mqcloudService.GetApplication(getApplicationOptionsModel)
				Expect(operationErr).To(BeNil())
				Expect(response).ToNot(BeNil())
				Expect(result).ToNot(BeNil())

				// Re-test the timeout error with retries disabled
				ctx, cancelFunc2 := context.WithTimeout(context.Background(), 80*time.Millisecond)
				defer cancelFunc2()
				_, _, operationErr = mqcloudService.GetApplicationWithContext(ctx, getApplicationOptionsModel)
				Expect(operationErr).ToNot(BeNil())
				Expect(operationErr.Error()).To(ContainSubstring("deadline exceeded"))
			})
			AfterEach(func() {
				testServer.Close()
			})
		})
		Context(`Using mock server endpoint`, func() {
			BeforeEach(func() {
				testServer = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
					defer GinkgoRecover()

					// Verify the contents of the request
					Expect(req.URL.EscapedPath()).To(Equal(getApplicationPath))
					Expect(req.Method).To(Equal("GET"))

					Expect(req.Header["Accept-Language"]).ToNot(BeNil())
					Expect(req.Header["Accept-Language"][0]).To(Equal(fmt.Sprintf("%v", "en-US,en;q=0.5")))
					// Set mock response
					res.Header().Set("Content-type", "application/json")
					res.WriteHeader(200)
					fmt.Fprintf(res, "%s", `{"id": "ID", "name": "Name", "iam_service_id": "IamServiceID", "create_api_key_uri": "CreateApiKeyURI", "href": "Href"}`)
				}))
			})
			It(`Invoke GetApplication successfully`, func() {
				mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1(&mqcloudv1.MqcloudV1Options{
					URL:            testServer.URL,
					Authenticator:  &core.NoAuthAuthenticator{},
					AcceptLanguage: core.StringPtr(acceptLanguage),
				})
				Expect(serviceErr).To(BeNil())
				Expect(mqcloudService).ToNot(BeNil())

				// Invoke operation with nil options model (negative test)
				result, response, operationErr := mqcloudService.GetApplication(nil)
				Expect(operationErr).NotTo(BeNil())
				Expect(response).To(BeNil())
				Expect(result).To(BeNil())

				// Construct an instance of the GetApplicationOptions model
				getApplicationOptionsModel := new(mqcloudv1.GetApplicationOptions)
				getApplicationOptionsModel.ServiceInstanceGuid = core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				getApplicationOptionsModel.ApplicationID = core.StringPtr("0123456789ABCDEF0123456789ABCDEF")
				getApplicationOptionsModel.Headers = map[string]string{"x-custom-header": "x-custom-value"}

				// Invoke operation with valid options model (positive test)
				result, response, operationErr = mqcloudService.GetApplication(getApplicationOptionsModel)
				Expect(operationErr).To(BeNil())
				Expect(response).ToNot(BeNil())
				Expect(result).ToNot(BeNil())

			})
			It(`Invoke GetApplication with error: Operation validation and request error`, func() {
				mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1(&mqcloudv1.MqcloudV1Options{
					URL:            testServer.URL,
					Authenticator:  &core.NoAuthAuthenticator{},
					AcceptLanguage: core.StringPtr(acceptLanguage),
				})
				Expect(serviceErr).To(BeNil())
				Expect(mqcloudService).ToNot(BeNil())

				// Construct an instance of the GetApplicationOptions model
				getApplicationOptionsModel := new(mqcloudv1.GetApplicationOptions)
				getApplicationOptionsModel.ServiceInstanceGuid = core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				getApplicationOptionsModel.ApplicationID = core.StringPtr("0123456789ABCDEF0123456789ABCDEF")
				getApplicationOptionsModel.Headers = map[string]string{"x-custom-header": "x-custom-value"}
				// Invoke operation with empty URL (negative test)
				err := mqcloudService.SetServiceURL("")
				Expect(err).To(BeNil())
				result, response, operationErr := mqcloudService.GetApplication(getApplicationOptionsModel)
				Expect(operationErr).ToNot(BeNil())
				Expect(operationErr.Error()).To(ContainSubstring(core.ERRORMSG_SERVICE_URL_MISSING))
				Expect(response).To(BeNil())
				Expect(result).To(BeNil())
				// Construct a second instance of the GetApplicationOptions model with no property values
				getApplicationOptionsModelNew := new(mqcloudv1.GetApplicationOptions)
				// Invoke operation with invalid model (negative test)
				result, response, operationErr = mqcloudService.GetApplication(getApplicationOptionsModelNew)
				Expect(operationErr).ToNot(BeNil())
				Expect(response).To(BeNil())
				Expect(result).To(BeNil())
			})
			AfterEach(func() {
				testServer.Close()
			})
		})
		Context(`Using mock server endpoint with missing response body`, func() {
			BeforeEach(func() {
				testServer = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
					defer GinkgoRecover()

					// Set success status code with no respoonse body
					res.WriteHeader(200)
				}))
			})
			It(`Invoke GetApplication successfully`, func() {
				mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1(&mqcloudv1.MqcloudV1Options{
					URL:            testServer.URL,
					Authenticator:  &core.NoAuthAuthenticator{},
					AcceptLanguage: core.StringPtr(acceptLanguage),
				})
				Expect(serviceErr).To(BeNil())
				Expect(mqcloudService).ToNot(BeNil())

				// Construct an instance of the GetApplicationOptions model
				getApplicationOptionsModel := new(mqcloudv1.GetApplicationOptions)
				getApplicationOptionsModel.ServiceInstanceGuid = core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				getApplicationOptionsModel.ApplicationID = core.StringPtr("0123456789ABCDEF0123456789ABCDEF")
				getApplicationOptionsModel.Headers = map[string]string{"x-custom-header": "x-custom-value"}

				// Invoke operation
				result, response, operationErr := mqcloudService.GetApplication(getApplicationOptionsModel)
				Expect(operationErr).To(BeNil())
				Expect(response).ToNot(BeNil())

				// Verify a nil result
				Expect(result).To(BeNil())
			})
			AfterEach(func() {
				testServer.Close()
			})
		})
	})
	Describe(`DeleteApplication(deleteApplicationOptions *DeleteApplicationOptions)`, func() {
		acceptLanguage := "en-US,en;q=0.5"
		deleteApplicationPath := "/v1/a2b4d4bc-dadb-4637-bcec-9b7d1e723af8/applications/0123456789ABCDEF0123456789ABCDEF"
		Context(`Using mock server endpoint`, func() {
			BeforeEach(func() {
				testServer = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
					defer GinkgoRecover()

					// Verify the contents of the request
					Expect(req.URL.EscapedPath()).To(Equal(deleteApplicationPath))
					Expect(req.Method).To(Equal("DELETE"))

					Expect(req.Header["Accept-Language"]).ToNot(BeNil())
					Expect(req.Header["Accept-Language"][0]).To(Equal(fmt.Sprintf("%v", "en-US,en;q=0.5")))
					res.WriteHeader(204)
				}))
			})
			It(`Invoke DeleteApplication successfully`, func() {
				mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1(&mqcloudv1.MqcloudV1Options{
					URL:            testServer.URL,
					Authenticator:  &core.NoAuthAuthenticator{},
					AcceptLanguage: core.StringPtr(acceptLanguage),
				})
				Expect(serviceErr).To(BeNil())
				Expect(mqcloudService).ToNot(BeNil())

				// Invoke operation with nil options model (negative test)
				response, operationErr := mqcloudService.DeleteApplication(nil)
				Expect(operationErr).NotTo(BeNil())
				Expect(response).To(BeNil())

				// Construct an instance of the DeleteApplicationOptions model
				deleteApplicationOptionsModel := new(mqcloudv1.DeleteApplicationOptions)
				deleteApplicationOptionsModel.ServiceInstanceGuid = core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				deleteApplicationOptionsModel.ApplicationID = core.StringPtr("0123456789ABCDEF0123456789ABCDEF")
				deleteApplicationOptionsModel.Headers = map[string]string{"x-custom-header": "x-custom-value"}

				// Invoke operation with valid options model (positive test)
				response, operationErr = mqcloudService.DeleteApplication(deleteApplicationOptionsModel)
				Expect(operationErr).To(BeNil())
				Expect(response).ToNot(BeNil())
			})
			It(`Invoke DeleteApplication with error: Operation validation and request error`, func() {
				mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1(&mqcloudv1.MqcloudV1Options{
					URL:            testServer.URL,
					Authenticator:  &core.NoAuthAuthenticator{},
					AcceptLanguage: core.StringPtr(acceptLanguage),
				})
				Expect(serviceErr).To(BeNil())
				Expect(mqcloudService).ToNot(BeNil())

				// Construct an instance of the DeleteApplicationOptions model
				deleteApplicationOptionsModel := new(mqcloudv1.DeleteApplicationOptions)
				deleteApplicationOptionsModel.ServiceInstanceGuid = core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				deleteApplicationOptionsModel.ApplicationID = core.StringPtr("0123456789ABCDEF0123456789ABCDEF")
				deleteApplicationOptionsModel.Headers = map[string]string{"x-custom-header": "x-custom-value"}
				// Invoke operation with empty URL (negative test)
				err := mqcloudService.SetServiceURL("")
				Expect(err).To(BeNil())
				response, operationErr := mqcloudService.DeleteApplication(deleteApplicationOptionsModel)
				Expect(operationErr).ToNot(BeNil())
				Expect(operationErr.Error()).To(ContainSubstring(core.ERRORMSG_SERVICE_URL_MISSING))
				Expect(response).To(BeNil())
				// Construct a second instance of the DeleteApplicationOptions model with no property values
				deleteApplicationOptionsModelNew := new(mqcloudv1.DeleteApplicationOptions)
				// Invoke operation with invalid model (negative test)
				response, operationErr = mqcloudService.DeleteApplication(deleteApplicationOptionsModelNew)
				Expect(operationErr).ToNot(BeNil())
				Expect(response).To(BeNil())
			})
			AfterEach(func() {
				testServer.Close()
			})
		})
	})
	Describe(`SetApplicationName(setApplicationNameOptions *SetApplicationNameOptions) - Operation response error`, func() {
		acceptLanguage := "en-US,en;q=0.5"
		setApplicationNamePath := "/v1/a2b4d4bc-dadb-4637-bcec-9b7d1e723af8/applications/0123456789ABCDEF0123456789ABCDEF"
		Context(`Using mock server endpoint with invalid JSON response`, func() {
			BeforeEach(func() {
				testServer = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
					defer GinkgoRecover()

					// Verify the contents of the request
					Expect(req.URL.EscapedPath()).To(Equal(setApplicationNamePath))
					Expect(req.Method).To(Equal("PATCH"))
					Expect(req.Header["Accept-Language"]).ToNot(BeNil())
					Expect(req.Header["Accept-Language"][0]).To(Equal(fmt.Sprintf("%v", "en-US,en;q=0.5")))
					res.Header().Set("Content-type", "application/json")
					res.WriteHeader(200)
					fmt.Fprint(res, `} this is not valid json {`)
				}))
			})
			It(`Invoke SetApplicationName with error: Operation response processing error`, func() {
				mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1(&mqcloudv1.MqcloudV1Options{
					URL:            testServer.URL,
					Authenticator:  &core.NoAuthAuthenticator{},
					AcceptLanguage: core.StringPtr(acceptLanguage),
				})
				Expect(serviceErr).To(BeNil())
				Expect(mqcloudService).ToNot(BeNil())

				// Construct an instance of the SetApplicationNameOptions model
				setApplicationNameOptionsModel := new(mqcloudv1.SetApplicationNameOptions)
				setApplicationNameOptionsModel.ServiceInstanceGuid = core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				setApplicationNameOptionsModel.ApplicationID = core.StringPtr("0123456789ABCDEF0123456789ABCDEF")
				setApplicationNameOptionsModel.Name = core.StringPtr("testString")
				setApplicationNameOptionsModel.Headers = map[string]string{"x-custom-header": "x-custom-value"}
				// Expect response parsing to fail since we are receiving a text/plain response
				result, response, operationErr := mqcloudService.SetApplicationName(setApplicationNameOptionsModel)
				Expect(operationErr).ToNot(BeNil())
				Expect(response).ToNot(BeNil())
				Expect(result).To(BeNil())

				// Enable retries and test again
				mqcloudService.EnableRetries(0, 0)
				result, response, operationErr = mqcloudService.SetApplicationName(setApplicationNameOptionsModel)
				Expect(operationErr).ToNot(BeNil())
				Expect(response).ToNot(BeNil())
				Expect(result).To(BeNil())
			})
			AfterEach(func() {
				testServer.Close()
			})
		})
	})
	Describe(`SetApplicationName(setApplicationNameOptions *SetApplicationNameOptions)`, func() {
		acceptLanguage := "en-US,en;q=0.5"
		setApplicationNamePath := "/v1/a2b4d4bc-dadb-4637-bcec-9b7d1e723af8/applications/0123456789ABCDEF0123456789ABCDEF"
		Context(`Using mock server endpoint with timeout`, func() {
			BeforeEach(func() {
				testServer = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
					defer GinkgoRecover()

					// Verify the contents of the request
					Expect(req.URL.EscapedPath()).To(Equal(setApplicationNamePath))
					Expect(req.Method).To(Equal("PATCH"))

					// For gzip-disabled operation, verify Content-Encoding is not set.
					Expect(req.Header.Get("Content-Encoding")).To(BeEmpty())

					// If there is a body, then make sure we can read it
					bodyBuf := new(bytes.Buffer)
					if req.Header.Get("Content-Encoding") == "gzip" {
						body, err := core.NewGzipDecompressionReader(req.Body)
						Expect(err).To(BeNil())
						_, err = bodyBuf.ReadFrom(body)
						Expect(err).To(BeNil())
					} else {
						_, err := bodyBuf.ReadFrom(req.Body)
						Expect(err).To(BeNil())
					}
					fmt.Fprintf(GinkgoWriter, "  Request body: %s", bodyBuf.String())

					Expect(req.Header["Accept-Language"]).ToNot(BeNil())
					Expect(req.Header["Accept-Language"][0]).To(Equal(fmt.Sprintf("%v", "en-US,en;q=0.5")))
					// Sleep a short time to support a timeout test
					time.Sleep(100 * time.Millisecond)

					// Set mock response
					res.Header().Set("Content-type", "application/json")
					res.WriteHeader(200)
					fmt.Fprintf(res, "%s", `{"id": "ID", "name": "Name", "iam_service_id": "IamServiceID", "create_api_key_uri": "CreateApiKeyURI", "href": "Href"}`)
				}))
			})
			It(`Invoke SetApplicationName successfully with retries`, func() {
				mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1(&mqcloudv1.MqcloudV1Options{
					URL:            testServer.URL,
					Authenticator:  &core.NoAuthAuthenticator{},
					AcceptLanguage: core.StringPtr(acceptLanguage),
				})
				Expect(serviceErr).To(BeNil())
				Expect(mqcloudService).ToNot(BeNil())
				mqcloudService.EnableRetries(0, 0)

				// Construct an instance of the SetApplicationNameOptions model
				setApplicationNameOptionsModel := new(mqcloudv1.SetApplicationNameOptions)
				setApplicationNameOptionsModel.ServiceInstanceGuid = core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				setApplicationNameOptionsModel.ApplicationID = core.StringPtr("0123456789ABCDEF0123456789ABCDEF")
				setApplicationNameOptionsModel.Name = core.StringPtr("testString")
				setApplicationNameOptionsModel.Headers = map[string]string{"x-custom-header": "x-custom-value"}

				// Invoke operation with a Context to test a timeout error
				ctx, cancelFunc := context.WithTimeout(context.Background(), 80*time.Millisecond)
				defer cancelFunc()
				_, _, operationErr := mqcloudService.SetApplicationNameWithContext(ctx, setApplicationNameOptionsModel)
				Expect(operationErr).ToNot(BeNil())
				Expect(operationErr.Error()).To(ContainSubstring("deadline exceeded"))

				// Disable retries and test again
				mqcloudService.DisableRetries()
				result, response, operationErr := mqcloudService.SetApplicationName(setApplicationNameOptionsModel)
				Expect(operationErr).To(BeNil())
				Expect(response).ToNot(BeNil())
				Expect(result).ToNot(BeNil())

				// Re-test the timeout error with retries disabled
				ctx, cancelFunc2 := context.WithTimeout(context.Background(), 80*time.Millisecond)
				defer cancelFunc2()
				_, _, operationErr = mqcloudService.SetApplicationNameWithContext(ctx, setApplicationNameOptionsModel)
				Expect(operationErr).ToNot(BeNil())
				Expect(operationErr.Error()).To(ContainSubstring("deadline exceeded"))
			})
			AfterEach(func() {
				testServer.Close()
			})
		})
		Context(`Using mock server endpoint`, func() {
			BeforeEach(func() {
				testServer = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
					defer GinkgoRecover()

					// Verify the contents of the request
					Expect(req.URL.EscapedPath()).To(Equal(setApplicationNamePath))
					Expect(req.Method).To(Equal("PATCH"))

					// For gzip-disabled operation, verify Content-Encoding is not set.
					Expect(req.Header.Get("Content-Encoding")).To(BeEmpty())

					// If there is a body, then make sure we can read it
					bodyBuf := new(bytes.Buffer)
					if req.Header.Get("Content-Encoding") == "gzip" {
						body, err := core.NewGzipDecompressionReader(req.Body)
						Expect(err).To(BeNil())
						_, err = bodyBuf.ReadFrom(body)
						Expect(err).To(BeNil())
					} else {
						_, err := bodyBuf.ReadFrom(req.Body)
						Expect(err).To(BeNil())
					}
					fmt.Fprintf(GinkgoWriter, "  Request body: %s", bodyBuf.String())

					Expect(req.Header["Accept-Language"]).ToNot(BeNil())
					Expect(req.Header["Accept-Language"][0]).To(Equal(fmt.Sprintf("%v", "en-US,en;q=0.5")))
					// Set mock response
					res.Header().Set("Content-type", "application/json")
					res.WriteHeader(200)
					fmt.Fprintf(res, "%s", `{"id": "ID", "name": "Name", "iam_service_id": "IamServiceID", "create_api_key_uri": "CreateApiKeyURI", "href": "Href"}`)
				}))
			})
			It(`Invoke SetApplicationName successfully`, func() {
				mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1(&mqcloudv1.MqcloudV1Options{
					URL:            testServer.URL,
					Authenticator:  &core.NoAuthAuthenticator{},
					AcceptLanguage: core.StringPtr(acceptLanguage),
				})
				Expect(serviceErr).To(BeNil())
				Expect(mqcloudService).ToNot(BeNil())

				// Invoke operation with nil options model (negative test)
				result, response, operationErr := mqcloudService.SetApplicationName(nil)
				Expect(operationErr).NotTo(BeNil())
				Expect(response).To(BeNil())
				Expect(result).To(BeNil())

				// Construct an instance of the SetApplicationNameOptions model
				setApplicationNameOptionsModel := new(mqcloudv1.SetApplicationNameOptions)
				setApplicationNameOptionsModel.ServiceInstanceGuid = core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				setApplicationNameOptionsModel.ApplicationID = core.StringPtr("0123456789ABCDEF0123456789ABCDEF")
				setApplicationNameOptionsModel.Name = core.StringPtr("testString")
				setApplicationNameOptionsModel.Headers = map[string]string{"x-custom-header": "x-custom-value"}

				// Invoke operation with valid options model (positive test)
				result, response, operationErr = mqcloudService.SetApplicationName(setApplicationNameOptionsModel)
				Expect(operationErr).To(BeNil())
				Expect(response).ToNot(BeNil())
				Expect(result).ToNot(BeNil())

			})
			It(`Invoke SetApplicationName with error: Operation validation and request error`, func() {
				mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1(&mqcloudv1.MqcloudV1Options{
					URL:            testServer.URL,
					Authenticator:  &core.NoAuthAuthenticator{},
					AcceptLanguage: core.StringPtr(acceptLanguage),
				})
				Expect(serviceErr).To(BeNil())
				Expect(mqcloudService).ToNot(BeNil())

				// Construct an instance of the SetApplicationNameOptions model
				setApplicationNameOptionsModel := new(mqcloudv1.SetApplicationNameOptions)
				setApplicationNameOptionsModel.ServiceInstanceGuid = core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				setApplicationNameOptionsModel.ApplicationID = core.StringPtr("0123456789ABCDEF0123456789ABCDEF")
				setApplicationNameOptionsModel.Name = core.StringPtr("testString")
				setApplicationNameOptionsModel.Headers = map[string]string{"x-custom-header": "x-custom-value"}
				// Invoke operation with empty URL (negative test)
				err := mqcloudService.SetServiceURL("")
				Expect(err).To(BeNil())
				result, response, operationErr := mqcloudService.SetApplicationName(setApplicationNameOptionsModel)
				Expect(operationErr).ToNot(BeNil())
				Expect(operationErr.Error()).To(ContainSubstring(core.ERRORMSG_SERVICE_URL_MISSING))
				Expect(response).To(BeNil())
				Expect(result).To(BeNil())
				// Construct a second instance of the SetApplicationNameOptions model with no property values
				setApplicationNameOptionsModelNew := new(mqcloudv1.SetApplicationNameOptions)
				// Invoke operation with invalid model (negative test)
				result, response, operationErr = mqcloudService.SetApplicationName(setApplicationNameOptionsModelNew)
				Expect(operationErr).ToNot(BeNil())
				Expect(response).To(BeNil())
				Expect(result).To(BeNil())
			})
			AfterEach(func() {
				testServer.Close()
			})
		})
		Context(`Using mock server endpoint with missing response body`, func() {
			BeforeEach(func() {
				testServer = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
					defer GinkgoRecover()

					// Set success status code with no respoonse body
					res.WriteHeader(200)
				}))
			})
			It(`Invoke SetApplicationName successfully`, func() {
				mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1(&mqcloudv1.MqcloudV1Options{
					URL:            testServer.URL,
					Authenticator:  &core.NoAuthAuthenticator{},
					AcceptLanguage: core.StringPtr(acceptLanguage),
				})
				Expect(serviceErr).To(BeNil())
				Expect(mqcloudService).ToNot(BeNil())

				// Construct an instance of the SetApplicationNameOptions model
				setApplicationNameOptionsModel := new(mqcloudv1.SetApplicationNameOptions)
				setApplicationNameOptionsModel.ServiceInstanceGuid = core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				setApplicationNameOptionsModel.ApplicationID = core.StringPtr("0123456789ABCDEF0123456789ABCDEF")
				setApplicationNameOptionsModel.Name = core.StringPtr("testString")
				setApplicationNameOptionsModel.Headers = map[string]string{"x-custom-header": "x-custom-value"}

				// Invoke operation
				result, response, operationErr := mqcloudService.SetApplicationName(setApplicationNameOptionsModel)
				Expect(operationErr).To(BeNil())
				Expect(response).ToNot(BeNil())

				// Verify a nil result
				Expect(result).To(BeNil())
			})
			AfterEach(func() {
				testServer.Close()
			})
		})
	})
	Describe(`CreateApplicationApikey(createApplicationApikeyOptions *CreateApplicationApikeyOptions) - Operation response error`, func() {
		acceptLanguage := "en-US,en;q=0.5"
		createApplicationApikeyPath := "/v1/a2b4d4bc-dadb-4637-bcec-9b7d1e723af8/applications/0123456789ABCDEF0123456789ABCDEF/api_key"
		Context(`Using mock server endpoint with invalid JSON response`, func() {
			BeforeEach(func() {
				testServer = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
					defer GinkgoRecover()

					// Verify the contents of the request
					Expect(req.URL.EscapedPath()).To(Equal(createApplicationApikeyPath))
					Expect(req.Method).To(Equal("POST"))
					Expect(req.Header["Accept-Language"]).ToNot(BeNil())
					Expect(req.Header["Accept-Language"][0]).To(Equal(fmt.Sprintf("%v", "en-US,en;q=0.5")))
					res.Header().Set("Content-type", "application/json")
					res.WriteHeader(201)
					fmt.Fprint(res, `} this is not valid json {`)
				}))
			})
			It(`Invoke CreateApplicationApikey with error: Operation response processing error`, func() {
				mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1(&mqcloudv1.MqcloudV1Options{
					URL:            testServer.URL,
					Authenticator:  &core.NoAuthAuthenticator{},
					AcceptLanguage: core.StringPtr(acceptLanguage),
				})
				Expect(serviceErr).To(BeNil())
				Expect(mqcloudService).ToNot(BeNil())

				// Construct an instance of the CreateApplicationApikeyOptions model
				createApplicationApikeyOptionsModel := new(mqcloudv1.CreateApplicationApikeyOptions)
				createApplicationApikeyOptionsModel.ServiceInstanceGuid = core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				createApplicationApikeyOptionsModel.ApplicationID = core.StringPtr("0123456789ABCDEF0123456789ABCDEF")
				createApplicationApikeyOptionsModel.Name = core.StringPtr("test-api-key")
				createApplicationApikeyOptionsModel.Headers = map[string]string{"x-custom-header": "x-custom-value"}
				// Expect response parsing to fail since we are receiving a text/plain response
				result, response, operationErr := mqcloudService.CreateApplicationApikey(createApplicationApikeyOptionsModel)
				Expect(operationErr).ToNot(BeNil())
				Expect(response).ToNot(BeNil())
				Expect(result).To(BeNil())

				// Enable retries and test again
				mqcloudService.EnableRetries(0, 0)
				result, response, operationErr = mqcloudService.CreateApplicationApikey(createApplicationApikeyOptionsModel)
				Expect(operationErr).ToNot(BeNil())
				Expect(response).ToNot(BeNil())
				Expect(result).To(BeNil())
			})
			AfterEach(func() {
				testServer.Close()
			})
		})
	})
	Describe(`CreateApplicationApikey(createApplicationApikeyOptions *CreateApplicationApikeyOptions)`, func() {
		acceptLanguage := "en-US,en;q=0.5"
		createApplicationApikeyPath := "/v1/a2b4d4bc-dadb-4637-bcec-9b7d1e723af8/applications/0123456789ABCDEF0123456789ABCDEF/api_key"
		Context(`Using mock server endpoint with timeout`, func() {
			BeforeEach(func() {
				testServer = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
					defer GinkgoRecover()

					// Verify the contents of the request
					Expect(req.URL.EscapedPath()).To(Equal(createApplicationApikeyPath))
					Expect(req.Method).To(Equal("POST"))

					// For gzip-disabled operation, verify Content-Encoding is not set.
					Expect(req.Header.Get("Content-Encoding")).To(BeEmpty())

					// If there is a body, then make sure we can read it
					bodyBuf := new(bytes.Buffer)
					if req.Header.Get("Content-Encoding") == "gzip" {
						body, err := core.NewGzipDecompressionReader(req.Body)
						Expect(err).To(BeNil())
						_, err = bodyBuf.ReadFrom(body)
						Expect(err).To(BeNil())
					} else {
						_, err := bodyBuf.ReadFrom(req.Body)
						Expect(err).To(BeNil())
					}
					fmt.Fprintf(GinkgoWriter, "  Request body: %s", bodyBuf.String())

					Expect(req.Header["Accept-Language"]).ToNot(BeNil())
					Expect(req.Header["Accept-Language"][0]).To(Equal(fmt.Sprintf("%v", "en-US,en;q=0.5")))
					// Sleep a short time to support a timeout test
					time.Sleep(100 * time.Millisecond)

					// Set mock response
					res.Header().Set("Content-type", "application/json")
					res.WriteHeader(201)
					fmt.Fprintf(res, "%s", `{"api_key_name": "ApiKeyName", "api_key_id": "ApiKeyID", "api_key": "ApiKey"}`)
				}))
			})
			It(`Invoke CreateApplicationApikey successfully with retries`, func() {
				mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1(&mqcloudv1.MqcloudV1Options{
					URL:            testServer.URL,
					Authenticator:  &core.NoAuthAuthenticator{},
					AcceptLanguage: core.StringPtr(acceptLanguage),
				})
				Expect(serviceErr).To(BeNil())
				Expect(mqcloudService).ToNot(BeNil())
				mqcloudService.EnableRetries(0, 0)

				// Construct an instance of the CreateApplicationApikeyOptions model
				createApplicationApikeyOptionsModel := new(mqcloudv1.CreateApplicationApikeyOptions)
				createApplicationApikeyOptionsModel.ServiceInstanceGuid = core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				createApplicationApikeyOptionsModel.ApplicationID = core.StringPtr("0123456789ABCDEF0123456789ABCDEF")
				createApplicationApikeyOptionsModel.Name = core.StringPtr("test-api-key")
				createApplicationApikeyOptionsModel.Headers = map[string]string{"x-custom-header": "x-custom-value"}

				// Invoke operation with a Context to test a timeout error
				ctx, cancelFunc := context.WithTimeout(context.Background(), 80*time.Millisecond)
				defer cancelFunc()
				_, _, operationErr := mqcloudService.CreateApplicationApikeyWithContext(ctx, createApplicationApikeyOptionsModel)
				Expect(operationErr).ToNot(BeNil())
				Expect(operationErr.Error()).To(ContainSubstring("deadline exceeded"))

				// Disable retries and test again
				mqcloudService.DisableRetries()
				result, response, operationErr := mqcloudService.CreateApplicationApikey(createApplicationApikeyOptionsModel)
				Expect(operationErr).To(BeNil())
				Expect(response).ToNot(BeNil())
				Expect(result).ToNot(BeNil())

				// Re-test the timeout error with retries disabled
				ctx, cancelFunc2 := context.WithTimeout(context.Background(), 80*time.Millisecond)
				defer cancelFunc2()
				_, _, operationErr = mqcloudService.CreateApplicationApikeyWithContext(ctx, createApplicationApikeyOptionsModel)
				Expect(operationErr).ToNot(BeNil())
				Expect(operationErr.Error()).To(ContainSubstring("deadline exceeded"))
			})
			AfterEach(func() {
				testServer.Close()
			})
		})
		Context(`Using mock server endpoint`, func() {
			BeforeEach(func() {
				testServer = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
					defer GinkgoRecover()

					// Verify the contents of the request
					Expect(req.URL.EscapedPath()).To(Equal(createApplicationApikeyPath))
					Expect(req.Method).To(Equal("POST"))

					// For gzip-disabled operation, verify Content-Encoding is not set.
					Expect(req.Header.Get("Content-Encoding")).To(BeEmpty())

					// If there is a body, then make sure we can read it
					bodyBuf := new(bytes.Buffer)
					if req.Header.Get("Content-Encoding") == "gzip" {
						body, err := core.NewGzipDecompressionReader(req.Body)
						Expect(err).To(BeNil())
						_, err = bodyBuf.ReadFrom(body)
						Expect(err).To(BeNil())
					} else {
						_, err := bodyBuf.ReadFrom(req.Body)
						Expect(err).To(BeNil())
					}
					fmt.Fprintf(GinkgoWriter, "  Request body: %s", bodyBuf.String())

					Expect(req.Header["Accept-Language"]).ToNot(BeNil())
					Expect(req.Header["Accept-Language"][0]).To(Equal(fmt.Sprintf("%v", "en-US,en;q=0.5")))
					// Set mock response
					res.Header().Set("Content-type", "application/json")
					res.WriteHeader(201)
					fmt.Fprintf(res, "%s", `{"api_key_name": "ApiKeyName", "api_key_id": "ApiKeyID", "api_key": "ApiKey"}`)
				}))
			})
			It(`Invoke CreateApplicationApikey successfully`, func() {
				mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1(&mqcloudv1.MqcloudV1Options{
					URL:            testServer.URL,
					Authenticator:  &core.NoAuthAuthenticator{},
					AcceptLanguage: core.StringPtr(acceptLanguage),
				})
				Expect(serviceErr).To(BeNil())
				Expect(mqcloudService).ToNot(BeNil())

				// Invoke operation with nil options model (negative test)
				result, response, operationErr := mqcloudService.CreateApplicationApikey(nil)
				Expect(operationErr).NotTo(BeNil())
				Expect(response).To(BeNil())
				Expect(result).To(BeNil())

				// Construct an instance of the CreateApplicationApikeyOptions model
				createApplicationApikeyOptionsModel := new(mqcloudv1.CreateApplicationApikeyOptions)
				createApplicationApikeyOptionsModel.ServiceInstanceGuid = core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				createApplicationApikeyOptionsModel.ApplicationID = core.StringPtr("0123456789ABCDEF0123456789ABCDEF")
				createApplicationApikeyOptionsModel.Name = core.StringPtr("test-api-key")
				createApplicationApikeyOptionsModel.Headers = map[string]string{"x-custom-header": "x-custom-value"}

				// Invoke operation with valid options model (positive test)
				result, response, operationErr = mqcloudService.CreateApplicationApikey(createApplicationApikeyOptionsModel)
				Expect(operationErr).To(BeNil())
				Expect(response).ToNot(BeNil())
				Expect(result).ToNot(BeNil())

			})
			It(`Invoke CreateApplicationApikey with error: Operation validation and request error`, func() {
				mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1(&mqcloudv1.MqcloudV1Options{
					URL:            testServer.URL,
					Authenticator:  &core.NoAuthAuthenticator{},
					AcceptLanguage: core.StringPtr(acceptLanguage),
				})
				Expect(serviceErr).To(BeNil())
				Expect(mqcloudService).ToNot(BeNil())

				// Construct an instance of the CreateApplicationApikeyOptions model
				createApplicationApikeyOptionsModel := new(mqcloudv1.CreateApplicationApikeyOptions)
				createApplicationApikeyOptionsModel.ServiceInstanceGuid = core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				createApplicationApikeyOptionsModel.ApplicationID = core.StringPtr("0123456789ABCDEF0123456789ABCDEF")
				createApplicationApikeyOptionsModel.Name = core.StringPtr("test-api-key")
				createApplicationApikeyOptionsModel.Headers = map[string]string{"x-custom-header": "x-custom-value"}
				// Invoke operation with empty URL (negative test)
				err := mqcloudService.SetServiceURL("")
				Expect(err).To(BeNil())
				result, response, operationErr := mqcloudService.CreateApplicationApikey(createApplicationApikeyOptionsModel)
				Expect(operationErr).ToNot(BeNil())
				Expect(operationErr.Error()).To(ContainSubstring(core.ERRORMSG_SERVICE_URL_MISSING))
				Expect(response).To(BeNil())
				Expect(result).To(BeNil())
				// Construct a second instance of the CreateApplicationApikeyOptions model with no property values
				createApplicationApikeyOptionsModelNew := new(mqcloudv1.CreateApplicationApikeyOptions)
				// Invoke operation with invalid model (negative test)
				result, response, operationErr = mqcloudService.CreateApplicationApikey(createApplicationApikeyOptionsModelNew)
				Expect(operationErr).ToNot(BeNil())
				Expect(response).To(BeNil())
				Expect(result).To(BeNil())
			})
			AfterEach(func() {
				testServer.Close()
			})
		})
		Context(`Using mock server endpoint with missing response body`, func() {
			BeforeEach(func() {
				testServer = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
					defer GinkgoRecover()

					// Set success status code with no respoonse body
					res.WriteHeader(201)
				}))
			})
			It(`Invoke CreateApplicationApikey successfully`, func() {
				mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1(&mqcloudv1.MqcloudV1Options{
					URL:            testServer.URL,
					Authenticator:  &core.NoAuthAuthenticator{},
					AcceptLanguage: core.StringPtr(acceptLanguage),
				})
				Expect(serviceErr).To(BeNil())
				Expect(mqcloudService).ToNot(BeNil())

				// Construct an instance of the CreateApplicationApikeyOptions model
				createApplicationApikeyOptionsModel := new(mqcloudv1.CreateApplicationApikeyOptions)
				createApplicationApikeyOptionsModel.ServiceInstanceGuid = core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				createApplicationApikeyOptionsModel.ApplicationID = core.StringPtr("0123456789ABCDEF0123456789ABCDEF")
				createApplicationApikeyOptionsModel.Name = core.StringPtr("test-api-key")
				createApplicationApikeyOptionsModel.Headers = map[string]string{"x-custom-header": "x-custom-value"}

				// Invoke operation
				result, response, operationErr := mqcloudService.CreateApplicationApikey(createApplicationApikeyOptionsModel)
				Expect(operationErr).To(BeNil())
				Expect(response).ToNot(BeNil())

				// Verify a nil result
				Expect(result).To(BeNil())
			})
			AfterEach(func() {
				testServer.Close()
			})
		})
	})
	Describe(`CreateTrustStorePemCertificate(createTrustStorePemCertificateOptions *CreateTrustStorePemCertificateOptions) - Operation response error`, func() {
		acceptLanguage := "en-US,en;q=0.5"
		createTrustStorePemCertificatePath := "/v1/a2b4d4bc-dadb-4637-bcec-9b7d1e723af8/queue_managers/b8e1aeda078009cf3db74e90d5d42328/certificates/trust_store"
		Context(`Using mock server endpoint with invalid JSON response`, func() {
			BeforeEach(func() {
				testServer = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
					defer GinkgoRecover()

					// Verify the contents of the request
					Expect(req.URL.EscapedPath()).To(Equal(createTrustStorePemCertificatePath))
					Expect(req.Method).To(Equal("POST"))
					Expect(req.Header["Accept-Language"]).ToNot(BeNil())
					Expect(req.Header["Accept-Language"][0]).To(Equal(fmt.Sprintf("%v", "en-US,en;q=0.5")))
					res.Header().Set("Content-type", "application/json")
					res.WriteHeader(201)
					fmt.Fprint(res, `} this is not valid json {`)
				}))
			})
			It(`Invoke CreateTrustStorePemCertificate with error: Operation response processing error`, func() {
				mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1(&mqcloudv1.MqcloudV1Options{
					URL:            testServer.URL,
					Authenticator:  &core.NoAuthAuthenticator{},
					AcceptLanguage: core.StringPtr(acceptLanguage),
				})
				Expect(serviceErr).To(BeNil())
				Expect(mqcloudService).ToNot(BeNil())

				// Construct an instance of the CreateTrustStorePemCertificateOptions model
				createTrustStorePemCertificateOptionsModel := new(mqcloudv1.CreateTrustStorePemCertificateOptions)
				createTrustStorePemCertificateOptionsModel.ServiceInstanceGuid = core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				createTrustStorePemCertificateOptionsModel.QueueManagerID = core.StringPtr("b8e1aeda078009cf3db74e90d5d42328")
				createTrustStorePemCertificateOptionsModel.Label = core.StringPtr("certlabel")
				createTrustStorePemCertificateOptionsModel.CertificateFile = CreateMockReader("This is a mock file.")
				createTrustStorePemCertificateOptionsModel.Headers = map[string]string{"x-custom-header": "x-custom-value"}
				// Expect response parsing to fail since we are receiving a text/plain response
				result, response, operationErr := mqcloudService.CreateTrustStorePemCertificate(createTrustStorePemCertificateOptionsModel)
				Expect(operationErr).ToNot(BeNil())
				Expect(response).ToNot(BeNil())
				Expect(result).To(BeNil())

				// Enable retries and test again
				mqcloudService.EnableRetries(0, 0)
				result, response, operationErr = mqcloudService.CreateTrustStorePemCertificate(createTrustStorePemCertificateOptionsModel)
				Expect(operationErr).ToNot(BeNil())
				Expect(response).ToNot(BeNil())
				Expect(result).To(BeNil())
			})
			AfterEach(func() {
				testServer.Close()
			})
		})
	})
	Describe(`CreateTrustStorePemCertificate(createTrustStorePemCertificateOptions *CreateTrustStorePemCertificateOptions)`, func() {
		acceptLanguage := "en-US,en;q=0.5"
		createTrustStorePemCertificatePath := "/v1/a2b4d4bc-dadb-4637-bcec-9b7d1e723af8/queue_managers/b8e1aeda078009cf3db74e90d5d42328/certificates/trust_store"
		Context(`Using mock server endpoint with timeout`, func() {
			BeforeEach(func() {
				testServer = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
					defer GinkgoRecover()

					// Verify the contents of the request
					Expect(req.URL.EscapedPath()).To(Equal(createTrustStorePemCertificatePath))
					Expect(req.Method).To(Equal("POST"))

					Expect(req.Header["Accept-Language"]).ToNot(BeNil())
					Expect(req.Header["Accept-Language"][0]).To(Equal(fmt.Sprintf("%v", "en-US,en;q=0.5")))
					// Sleep a short time to support a timeout test
					time.Sleep(100 * time.Millisecond)

					// Set mock response
					res.Header().Set("Content-type", "application/json")
					res.WriteHeader(201)
					fmt.Fprintf(res, "%s", `{"id": "ID", "label": "Label", "certificate_type": "trust_store", "fingerprint_sha256": "FingerprintSha256", "subject_dn": "SubjectDn", "subject_cn": "SubjectCn", "issuer_dn": "IssuerDn", "issuer_cn": "IssuerCn", "issued": "2019-01-01T12:00:00.000Z", "expiry": "2019-01-01T12:00:00.000Z", "trusted": false, "href": "Href"}`)
				}))
			})
			It(`Invoke CreateTrustStorePemCertificate successfully with retries`, func() {
				mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1(&mqcloudv1.MqcloudV1Options{
					URL:            testServer.URL,
					Authenticator:  &core.NoAuthAuthenticator{},
					AcceptLanguage: core.StringPtr(acceptLanguage),
				})
				Expect(serviceErr).To(BeNil())
				Expect(mqcloudService).ToNot(BeNil())
				mqcloudService.EnableRetries(0, 0)

				// Construct an instance of the CreateTrustStorePemCertificateOptions model
				createTrustStorePemCertificateOptionsModel := new(mqcloudv1.CreateTrustStorePemCertificateOptions)
				createTrustStorePemCertificateOptionsModel.ServiceInstanceGuid = core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				createTrustStorePemCertificateOptionsModel.QueueManagerID = core.StringPtr("b8e1aeda078009cf3db74e90d5d42328")
				createTrustStorePemCertificateOptionsModel.Label = core.StringPtr("certlabel")
				createTrustStorePemCertificateOptionsModel.CertificateFile = CreateMockReader("This is a mock file.")
				createTrustStorePemCertificateOptionsModel.Headers = map[string]string{"x-custom-header": "x-custom-value"}

				// Invoke operation with a Context to test a timeout error
				ctx, cancelFunc := context.WithTimeout(context.Background(), 80*time.Millisecond)
				defer cancelFunc()
				_, _, operationErr := mqcloudService.CreateTrustStorePemCertificateWithContext(ctx, createTrustStorePemCertificateOptionsModel)
				Expect(operationErr).ToNot(BeNil())
				Expect(operationErr.Error()).To(ContainSubstring("deadline exceeded"))

				// Disable retries and test again
				mqcloudService.DisableRetries()
				result, response, operationErr := mqcloudService.CreateTrustStorePemCertificate(createTrustStorePemCertificateOptionsModel)
				Expect(operationErr).To(BeNil())
				Expect(response).ToNot(BeNil())
				Expect(result).ToNot(BeNil())

				// Re-test the timeout error with retries disabled
				ctx, cancelFunc2 := context.WithTimeout(context.Background(), 80*time.Millisecond)
				defer cancelFunc2()
				_, _, operationErr = mqcloudService.CreateTrustStorePemCertificateWithContext(ctx, createTrustStorePemCertificateOptionsModel)
				Expect(operationErr).ToNot(BeNil())
				Expect(operationErr.Error()).To(ContainSubstring("deadline exceeded"))
			})
			AfterEach(func() {
				testServer.Close()
			})
		})
		Context(`Using mock server endpoint`, func() {
			BeforeEach(func() {
				testServer = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
					defer GinkgoRecover()

					// Verify the contents of the request
					Expect(req.URL.EscapedPath()).To(Equal(createTrustStorePemCertificatePath))
					Expect(req.Method).To(Equal("POST"))

					Expect(req.Header["Accept-Language"]).ToNot(BeNil())
					Expect(req.Header["Accept-Language"][0]).To(Equal(fmt.Sprintf("%v", "en-US,en;q=0.5")))
					// Set mock response
					res.Header().Set("Content-type", "application/json")
					res.WriteHeader(201)
					fmt.Fprintf(res, "%s", `{"id": "ID", "label": "Label", "certificate_type": "trust_store", "fingerprint_sha256": "FingerprintSha256", "subject_dn": "SubjectDn", "subject_cn": "SubjectCn", "issuer_dn": "IssuerDn", "issuer_cn": "IssuerCn", "issued": "2019-01-01T12:00:00.000Z", "expiry": "2019-01-01T12:00:00.000Z", "trusted": false, "href": "Href"}`)
				}))
			})
			It(`Invoke CreateTrustStorePemCertificate successfully`, func() {
				mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1(&mqcloudv1.MqcloudV1Options{
					URL:            testServer.URL,
					Authenticator:  &core.NoAuthAuthenticator{},
					AcceptLanguage: core.StringPtr(acceptLanguage),
				})
				Expect(serviceErr).To(BeNil())
				Expect(mqcloudService).ToNot(BeNil())

				// Invoke operation with nil options model (negative test)
				result, response, operationErr := mqcloudService.CreateTrustStorePemCertificate(nil)
				Expect(operationErr).NotTo(BeNil())
				Expect(response).To(BeNil())
				Expect(result).To(BeNil())

				// Construct an instance of the CreateTrustStorePemCertificateOptions model
				createTrustStorePemCertificateOptionsModel := new(mqcloudv1.CreateTrustStorePemCertificateOptions)
				createTrustStorePemCertificateOptionsModel.ServiceInstanceGuid = core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				createTrustStorePemCertificateOptionsModel.QueueManagerID = core.StringPtr("b8e1aeda078009cf3db74e90d5d42328")
				createTrustStorePemCertificateOptionsModel.Label = core.StringPtr("certlabel")
				createTrustStorePemCertificateOptionsModel.CertificateFile = CreateMockReader("This is a mock file.")
				createTrustStorePemCertificateOptionsModel.Headers = map[string]string{"x-custom-header": "x-custom-value"}

				// Invoke operation with valid options model (positive test)
				result, response, operationErr = mqcloudService.CreateTrustStorePemCertificate(createTrustStorePemCertificateOptionsModel)
				Expect(operationErr).To(BeNil())
				Expect(response).ToNot(BeNil())
				Expect(result).ToNot(BeNil())

			})
			It(`Invoke CreateTrustStorePemCertificate with error: Operation validation and request error`, func() {
				mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1(&mqcloudv1.MqcloudV1Options{
					URL:            testServer.URL,
					Authenticator:  &core.NoAuthAuthenticator{},
					AcceptLanguage: core.StringPtr(acceptLanguage),
				})
				Expect(serviceErr).To(BeNil())
				Expect(mqcloudService).ToNot(BeNil())

				// Construct an instance of the CreateTrustStorePemCertificateOptions model
				createTrustStorePemCertificateOptionsModel := new(mqcloudv1.CreateTrustStorePemCertificateOptions)
				createTrustStorePemCertificateOptionsModel.ServiceInstanceGuid = core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				createTrustStorePemCertificateOptionsModel.QueueManagerID = core.StringPtr("b8e1aeda078009cf3db74e90d5d42328")
				createTrustStorePemCertificateOptionsModel.Label = core.StringPtr("certlabel")
				createTrustStorePemCertificateOptionsModel.CertificateFile = CreateMockReader("This is a mock file.")
				createTrustStorePemCertificateOptionsModel.Headers = map[string]string{"x-custom-header": "x-custom-value"}
				// Invoke operation with empty URL (negative test)
				err := mqcloudService.SetServiceURL("")
				Expect(err).To(BeNil())
				result, response, operationErr := mqcloudService.CreateTrustStorePemCertificate(createTrustStorePemCertificateOptionsModel)
				Expect(operationErr).ToNot(BeNil())
				Expect(operationErr.Error()).To(ContainSubstring(core.ERRORMSG_SERVICE_URL_MISSING))
				Expect(response).To(BeNil())
				Expect(result).To(BeNil())
				// Construct a second instance of the CreateTrustStorePemCertificateOptions model with no property values
				createTrustStorePemCertificateOptionsModelNew := new(mqcloudv1.CreateTrustStorePemCertificateOptions)
				// Invoke operation with invalid model (negative test)
				result, response, operationErr = mqcloudService.CreateTrustStorePemCertificate(createTrustStorePemCertificateOptionsModelNew)
				Expect(operationErr).ToNot(BeNil())
				Expect(response).To(BeNil())
				Expect(result).To(BeNil())
			})
			AfterEach(func() {
				testServer.Close()
			})
		})
		Context(`Using mock server endpoint with missing response body`, func() {
			BeforeEach(func() {
				testServer = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
					defer GinkgoRecover()

					// Set success status code with no respoonse body
					res.WriteHeader(201)
				}))
			})
			It(`Invoke CreateTrustStorePemCertificate successfully`, func() {
				mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1(&mqcloudv1.MqcloudV1Options{
					URL:            testServer.URL,
					Authenticator:  &core.NoAuthAuthenticator{},
					AcceptLanguage: core.StringPtr(acceptLanguage),
				})
				Expect(serviceErr).To(BeNil())
				Expect(mqcloudService).ToNot(BeNil())

				// Construct an instance of the CreateTrustStorePemCertificateOptions model
				createTrustStorePemCertificateOptionsModel := new(mqcloudv1.CreateTrustStorePemCertificateOptions)
				createTrustStorePemCertificateOptionsModel.ServiceInstanceGuid = core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				createTrustStorePemCertificateOptionsModel.QueueManagerID = core.StringPtr("b8e1aeda078009cf3db74e90d5d42328")
				createTrustStorePemCertificateOptionsModel.Label = core.StringPtr("certlabel")
				createTrustStorePemCertificateOptionsModel.CertificateFile = CreateMockReader("This is a mock file.")
				createTrustStorePemCertificateOptionsModel.Headers = map[string]string{"x-custom-header": "x-custom-value"}

				// Invoke operation
				result, response, operationErr := mqcloudService.CreateTrustStorePemCertificate(createTrustStorePemCertificateOptionsModel)
				Expect(operationErr).To(BeNil())
				Expect(response).ToNot(BeNil())

				// Verify a nil result
				Expect(result).To(BeNil())
			})
			AfterEach(func() {
				testServer.Close()
			})
		})
	})
	Describe(`ListTrustStoreCertificates(listTrustStoreCertificatesOptions *ListTrustStoreCertificatesOptions) - Operation response error`, func() {
		acceptLanguage := "en-US,en;q=0.5"
		listTrustStoreCertificatesPath := "/v1/a2b4d4bc-dadb-4637-bcec-9b7d1e723af8/queue_managers/b8e1aeda078009cf3db74e90d5d42328/certificates/trust_store"
		Context(`Using mock server endpoint with invalid JSON response`, func() {
			BeforeEach(func() {
				testServer = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
					defer GinkgoRecover()

					// Verify the contents of the request
					Expect(req.URL.EscapedPath()).To(Equal(listTrustStoreCertificatesPath))
					Expect(req.Method).To(Equal("GET"))
					Expect(req.Header["Accept-Language"]).ToNot(BeNil())
					Expect(req.Header["Accept-Language"][0]).To(Equal(fmt.Sprintf("%v", "en-US,en;q=0.5")))
					res.Header().Set("Content-type", "application/json")
					res.WriteHeader(200)
					fmt.Fprint(res, `} this is not valid json {`)
				}))
			})
			It(`Invoke ListTrustStoreCertificates with error: Operation response processing error`, func() {
				mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1(&mqcloudv1.MqcloudV1Options{
					URL:            testServer.URL,
					Authenticator:  &core.NoAuthAuthenticator{},
					AcceptLanguage: core.StringPtr(acceptLanguage),
				})
				Expect(serviceErr).To(BeNil())
				Expect(mqcloudService).ToNot(BeNil())

				// Construct an instance of the ListTrustStoreCertificatesOptions model
				listTrustStoreCertificatesOptionsModel := new(mqcloudv1.ListTrustStoreCertificatesOptions)
				listTrustStoreCertificatesOptionsModel.ServiceInstanceGuid = core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				listTrustStoreCertificatesOptionsModel.QueueManagerID = core.StringPtr("b8e1aeda078009cf3db74e90d5d42328")
				listTrustStoreCertificatesOptionsModel.Headers = map[string]string{"x-custom-header": "x-custom-value"}
				// Expect response parsing to fail since we are receiving a text/plain response
				result, response, operationErr := mqcloudService.ListTrustStoreCertificates(listTrustStoreCertificatesOptionsModel)
				Expect(operationErr).ToNot(BeNil())
				Expect(response).ToNot(BeNil())
				Expect(result).To(BeNil())

				// Enable retries and test again
				mqcloudService.EnableRetries(0, 0)
				result, response, operationErr = mqcloudService.ListTrustStoreCertificates(listTrustStoreCertificatesOptionsModel)
				Expect(operationErr).ToNot(BeNil())
				Expect(response).ToNot(BeNil())
				Expect(result).To(BeNil())
			})
			AfterEach(func() {
				testServer.Close()
			})
		})
	})
	Describe(`ListTrustStoreCertificates(listTrustStoreCertificatesOptions *ListTrustStoreCertificatesOptions)`, func() {
		acceptLanguage := "en-US,en;q=0.5"
		listTrustStoreCertificatesPath := "/v1/a2b4d4bc-dadb-4637-bcec-9b7d1e723af8/queue_managers/b8e1aeda078009cf3db74e90d5d42328/certificates/trust_store"
		Context(`Using mock server endpoint with timeout`, func() {
			BeforeEach(func() {
				testServer = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
					defer GinkgoRecover()

					// Verify the contents of the request
					Expect(req.URL.EscapedPath()).To(Equal(listTrustStoreCertificatesPath))
					Expect(req.Method).To(Equal("GET"))

					Expect(req.Header["Accept-Language"]).ToNot(BeNil())
					Expect(req.Header["Accept-Language"][0]).To(Equal(fmt.Sprintf("%v", "en-US,en;q=0.5")))
					// Sleep a short time to support a timeout test
					time.Sleep(100 * time.Millisecond)

					// Set mock response
					res.Header().Set("Content-type", "application/json")
					res.WriteHeader(200)
					fmt.Fprintf(res, "%s", `{"total_count": 1, "trust_store": [{"id": "ID", "label": "Label", "certificate_type": "trust_store", "fingerprint_sha256": "FingerprintSha256", "subject_dn": "SubjectDn", "subject_cn": "SubjectCn", "issuer_dn": "IssuerDn", "issuer_cn": "IssuerCn", "issued": "2019-01-01T12:00:00.000Z", "expiry": "2019-01-01T12:00:00.000Z", "trusted": false, "href": "Href"}]}`)
				}))
			})
			It(`Invoke ListTrustStoreCertificates successfully with retries`, func() {
				mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1(&mqcloudv1.MqcloudV1Options{
					URL:            testServer.URL,
					Authenticator:  &core.NoAuthAuthenticator{},
					AcceptLanguage: core.StringPtr(acceptLanguage),
				})
				Expect(serviceErr).To(BeNil())
				Expect(mqcloudService).ToNot(BeNil())
				mqcloudService.EnableRetries(0, 0)

				// Construct an instance of the ListTrustStoreCertificatesOptions model
				listTrustStoreCertificatesOptionsModel := new(mqcloudv1.ListTrustStoreCertificatesOptions)
				listTrustStoreCertificatesOptionsModel.ServiceInstanceGuid = core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				listTrustStoreCertificatesOptionsModel.QueueManagerID = core.StringPtr("b8e1aeda078009cf3db74e90d5d42328")
				listTrustStoreCertificatesOptionsModel.Headers = map[string]string{"x-custom-header": "x-custom-value"}

				// Invoke operation with a Context to test a timeout error
				ctx, cancelFunc := context.WithTimeout(context.Background(), 80*time.Millisecond)
				defer cancelFunc()
				_, _, operationErr := mqcloudService.ListTrustStoreCertificatesWithContext(ctx, listTrustStoreCertificatesOptionsModel)
				Expect(operationErr).ToNot(BeNil())
				Expect(operationErr.Error()).To(ContainSubstring("deadline exceeded"))

				// Disable retries and test again
				mqcloudService.DisableRetries()
				result, response, operationErr := mqcloudService.ListTrustStoreCertificates(listTrustStoreCertificatesOptionsModel)
				Expect(operationErr).To(BeNil())
				Expect(response).ToNot(BeNil())
				Expect(result).ToNot(BeNil())

				// Re-test the timeout error with retries disabled
				ctx, cancelFunc2 := context.WithTimeout(context.Background(), 80*time.Millisecond)
				defer cancelFunc2()
				_, _, operationErr = mqcloudService.ListTrustStoreCertificatesWithContext(ctx, listTrustStoreCertificatesOptionsModel)
				Expect(operationErr).ToNot(BeNil())
				Expect(operationErr.Error()).To(ContainSubstring("deadline exceeded"))
			})
			AfterEach(func() {
				testServer.Close()
			})
		})
		Context(`Using mock server endpoint`, func() {
			BeforeEach(func() {
				testServer = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
					defer GinkgoRecover()

					// Verify the contents of the request
					Expect(req.URL.EscapedPath()).To(Equal(listTrustStoreCertificatesPath))
					Expect(req.Method).To(Equal("GET"))

					Expect(req.Header["Accept-Language"]).ToNot(BeNil())
					Expect(req.Header["Accept-Language"][0]).To(Equal(fmt.Sprintf("%v", "en-US,en;q=0.5")))
					// Set mock response
					res.Header().Set("Content-type", "application/json")
					res.WriteHeader(200)
					fmt.Fprintf(res, "%s", `{"total_count": 1, "trust_store": [{"id": "ID", "label": "Label", "certificate_type": "trust_store", "fingerprint_sha256": "FingerprintSha256", "subject_dn": "SubjectDn", "subject_cn": "SubjectCn", "issuer_dn": "IssuerDn", "issuer_cn": "IssuerCn", "issued": "2019-01-01T12:00:00.000Z", "expiry": "2019-01-01T12:00:00.000Z", "trusted": false, "href": "Href"}]}`)
				}))
			})
			It(`Invoke ListTrustStoreCertificates successfully`, func() {
				mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1(&mqcloudv1.MqcloudV1Options{
					URL:            testServer.URL,
					Authenticator:  &core.NoAuthAuthenticator{},
					AcceptLanguage: core.StringPtr(acceptLanguage),
				})
				Expect(serviceErr).To(BeNil())
				Expect(mqcloudService).ToNot(BeNil())

				// Invoke operation with nil options model (negative test)
				result, response, operationErr := mqcloudService.ListTrustStoreCertificates(nil)
				Expect(operationErr).NotTo(BeNil())
				Expect(response).To(BeNil())
				Expect(result).To(BeNil())

				// Construct an instance of the ListTrustStoreCertificatesOptions model
				listTrustStoreCertificatesOptionsModel := new(mqcloudv1.ListTrustStoreCertificatesOptions)
				listTrustStoreCertificatesOptionsModel.ServiceInstanceGuid = core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				listTrustStoreCertificatesOptionsModel.QueueManagerID = core.StringPtr("b8e1aeda078009cf3db74e90d5d42328")
				listTrustStoreCertificatesOptionsModel.Headers = map[string]string{"x-custom-header": "x-custom-value"}

				// Invoke operation with valid options model (positive test)
				result, response, operationErr = mqcloudService.ListTrustStoreCertificates(listTrustStoreCertificatesOptionsModel)
				Expect(operationErr).To(BeNil())
				Expect(response).ToNot(BeNil())
				Expect(result).ToNot(BeNil())

			})
			It(`Invoke ListTrustStoreCertificates with error: Operation validation and request error`, func() {
				mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1(&mqcloudv1.MqcloudV1Options{
					URL:            testServer.URL,
					Authenticator:  &core.NoAuthAuthenticator{},
					AcceptLanguage: core.StringPtr(acceptLanguage),
				})
				Expect(serviceErr).To(BeNil())
				Expect(mqcloudService).ToNot(BeNil())

				// Construct an instance of the ListTrustStoreCertificatesOptions model
				listTrustStoreCertificatesOptionsModel := new(mqcloudv1.ListTrustStoreCertificatesOptions)
				listTrustStoreCertificatesOptionsModel.ServiceInstanceGuid = core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				listTrustStoreCertificatesOptionsModel.QueueManagerID = core.StringPtr("b8e1aeda078009cf3db74e90d5d42328")
				listTrustStoreCertificatesOptionsModel.Headers = map[string]string{"x-custom-header": "x-custom-value"}
				// Invoke operation with empty URL (negative test)
				err := mqcloudService.SetServiceURL("")
				Expect(err).To(BeNil())
				result, response, operationErr := mqcloudService.ListTrustStoreCertificates(listTrustStoreCertificatesOptionsModel)
				Expect(operationErr).ToNot(BeNil())
				Expect(operationErr.Error()).To(ContainSubstring(core.ERRORMSG_SERVICE_URL_MISSING))
				Expect(response).To(BeNil())
				Expect(result).To(BeNil())
				// Construct a second instance of the ListTrustStoreCertificatesOptions model with no property values
				listTrustStoreCertificatesOptionsModelNew := new(mqcloudv1.ListTrustStoreCertificatesOptions)
				// Invoke operation with invalid model (negative test)
				result, response, operationErr = mqcloudService.ListTrustStoreCertificates(listTrustStoreCertificatesOptionsModelNew)
				Expect(operationErr).ToNot(BeNil())
				Expect(response).To(BeNil())
				Expect(result).To(BeNil())
			})
			AfterEach(func() {
				testServer.Close()
			})
		})
		Context(`Using mock server endpoint with missing response body`, func() {
			BeforeEach(func() {
				testServer = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
					defer GinkgoRecover()

					// Set success status code with no respoonse body
					res.WriteHeader(200)
				}))
			})
			It(`Invoke ListTrustStoreCertificates successfully`, func() {
				mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1(&mqcloudv1.MqcloudV1Options{
					URL:            testServer.URL,
					Authenticator:  &core.NoAuthAuthenticator{},
					AcceptLanguage: core.StringPtr(acceptLanguage),
				})
				Expect(serviceErr).To(BeNil())
				Expect(mqcloudService).ToNot(BeNil())

				// Construct an instance of the ListTrustStoreCertificatesOptions model
				listTrustStoreCertificatesOptionsModel := new(mqcloudv1.ListTrustStoreCertificatesOptions)
				listTrustStoreCertificatesOptionsModel.ServiceInstanceGuid = core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				listTrustStoreCertificatesOptionsModel.QueueManagerID = core.StringPtr("b8e1aeda078009cf3db74e90d5d42328")
				listTrustStoreCertificatesOptionsModel.Headers = map[string]string{"x-custom-header": "x-custom-value"}

				// Invoke operation
				result, response, operationErr := mqcloudService.ListTrustStoreCertificates(listTrustStoreCertificatesOptionsModel)
				Expect(operationErr).To(BeNil())
				Expect(response).ToNot(BeNil())

				// Verify a nil result
				Expect(result).To(BeNil())
			})
			AfterEach(func() {
				testServer.Close()
			})
		})
	})
	Describe(`GetTrustStoreCertificate(getTrustStoreCertificateOptions *GetTrustStoreCertificateOptions) - Operation response error`, func() {
		acceptLanguage := "en-US,en;q=0.5"
		getTrustStoreCertificatePath := "/v1/a2b4d4bc-dadb-4637-bcec-9b7d1e723af8/queue_managers/b8e1aeda078009cf3db74e90d5d42328/certificates/trust_store/9b7d1e723af8233"
		Context(`Using mock server endpoint with invalid JSON response`, func() {
			BeforeEach(func() {
				testServer = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
					defer GinkgoRecover()

					// Verify the contents of the request
					Expect(req.URL.EscapedPath()).To(Equal(getTrustStoreCertificatePath))
					Expect(req.Method).To(Equal("GET"))
					Expect(req.Header["Accept-Language"]).ToNot(BeNil())
					Expect(req.Header["Accept-Language"][0]).To(Equal(fmt.Sprintf("%v", "en-US,en;q=0.5")))
					res.Header().Set("Content-type", "application/json")
					res.WriteHeader(200)
					fmt.Fprint(res, `} this is not valid json {`)
				}))
			})
			It(`Invoke GetTrustStoreCertificate with error: Operation response processing error`, func() {
				mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1(&mqcloudv1.MqcloudV1Options{
					URL:            testServer.URL,
					Authenticator:  &core.NoAuthAuthenticator{},
					AcceptLanguage: core.StringPtr(acceptLanguage),
				})
				Expect(serviceErr).To(BeNil())
				Expect(mqcloudService).ToNot(BeNil())

				// Construct an instance of the GetTrustStoreCertificateOptions model
				getTrustStoreCertificateOptionsModel := new(mqcloudv1.GetTrustStoreCertificateOptions)
				getTrustStoreCertificateOptionsModel.ServiceInstanceGuid = core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				getTrustStoreCertificateOptionsModel.QueueManagerID = core.StringPtr("b8e1aeda078009cf3db74e90d5d42328")
				getTrustStoreCertificateOptionsModel.CertificateID = core.StringPtr("9b7d1e723af8233")
				getTrustStoreCertificateOptionsModel.Headers = map[string]string{"x-custom-header": "x-custom-value"}
				// Expect response parsing to fail since we are receiving a text/plain response
				result, response, operationErr := mqcloudService.GetTrustStoreCertificate(getTrustStoreCertificateOptionsModel)
				Expect(operationErr).ToNot(BeNil())
				Expect(response).ToNot(BeNil())
				Expect(result).To(BeNil())

				// Enable retries and test again
				mqcloudService.EnableRetries(0, 0)
				result, response, operationErr = mqcloudService.GetTrustStoreCertificate(getTrustStoreCertificateOptionsModel)
				Expect(operationErr).ToNot(BeNil())
				Expect(response).ToNot(BeNil())
				Expect(result).To(BeNil())
			})
			AfterEach(func() {
				testServer.Close()
			})
		})
	})
	Describe(`GetTrustStoreCertificate(getTrustStoreCertificateOptions *GetTrustStoreCertificateOptions)`, func() {
		acceptLanguage := "en-US,en;q=0.5"
		getTrustStoreCertificatePath := "/v1/a2b4d4bc-dadb-4637-bcec-9b7d1e723af8/queue_managers/b8e1aeda078009cf3db74e90d5d42328/certificates/trust_store/9b7d1e723af8233"
		Context(`Using mock server endpoint with timeout`, func() {
			BeforeEach(func() {
				testServer = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
					defer GinkgoRecover()

					// Verify the contents of the request
					Expect(req.URL.EscapedPath()).To(Equal(getTrustStoreCertificatePath))
					Expect(req.Method).To(Equal("GET"))

					Expect(req.Header["Accept-Language"]).ToNot(BeNil())
					Expect(req.Header["Accept-Language"][0]).To(Equal(fmt.Sprintf("%v", "en-US,en;q=0.5")))
					// Sleep a short time to support a timeout test
					time.Sleep(100 * time.Millisecond)

					// Set mock response
					res.Header().Set("Content-type", "application/json")
					res.WriteHeader(200)
					fmt.Fprintf(res, "%s", `{"id": "ID", "label": "Label", "certificate_type": "trust_store", "fingerprint_sha256": "FingerprintSha256", "subject_dn": "SubjectDn", "subject_cn": "SubjectCn", "issuer_dn": "IssuerDn", "issuer_cn": "IssuerCn", "issued": "2019-01-01T12:00:00.000Z", "expiry": "2019-01-01T12:00:00.000Z", "trusted": false, "href": "Href"}`)
				}))
			})
			It(`Invoke GetTrustStoreCertificate successfully with retries`, func() {
				mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1(&mqcloudv1.MqcloudV1Options{
					URL:            testServer.URL,
					Authenticator:  &core.NoAuthAuthenticator{},
					AcceptLanguage: core.StringPtr(acceptLanguage),
				})
				Expect(serviceErr).To(BeNil())
				Expect(mqcloudService).ToNot(BeNil())
				mqcloudService.EnableRetries(0, 0)

				// Construct an instance of the GetTrustStoreCertificateOptions model
				getTrustStoreCertificateOptionsModel := new(mqcloudv1.GetTrustStoreCertificateOptions)
				getTrustStoreCertificateOptionsModel.ServiceInstanceGuid = core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				getTrustStoreCertificateOptionsModel.QueueManagerID = core.StringPtr("b8e1aeda078009cf3db74e90d5d42328")
				getTrustStoreCertificateOptionsModel.CertificateID = core.StringPtr("9b7d1e723af8233")
				getTrustStoreCertificateOptionsModel.Headers = map[string]string{"x-custom-header": "x-custom-value"}

				// Invoke operation with a Context to test a timeout error
				ctx, cancelFunc := context.WithTimeout(context.Background(), 80*time.Millisecond)
				defer cancelFunc()
				_, _, operationErr := mqcloudService.GetTrustStoreCertificateWithContext(ctx, getTrustStoreCertificateOptionsModel)
				Expect(operationErr).ToNot(BeNil())
				Expect(operationErr.Error()).To(ContainSubstring("deadline exceeded"))

				// Disable retries and test again
				mqcloudService.DisableRetries()
				result, response, operationErr := mqcloudService.GetTrustStoreCertificate(getTrustStoreCertificateOptionsModel)
				Expect(operationErr).To(BeNil())
				Expect(response).ToNot(BeNil())
				Expect(result).ToNot(BeNil())

				// Re-test the timeout error with retries disabled
				ctx, cancelFunc2 := context.WithTimeout(context.Background(), 80*time.Millisecond)
				defer cancelFunc2()
				_, _, operationErr = mqcloudService.GetTrustStoreCertificateWithContext(ctx, getTrustStoreCertificateOptionsModel)
				Expect(operationErr).ToNot(BeNil())
				Expect(operationErr.Error()).To(ContainSubstring("deadline exceeded"))
			})
			AfterEach(func() {
				testServer.Close()
			})
		})
		Context(`Using mock server endpoint`, func() {
			BeforeEach(func() {
				testServer = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
					defer GinkgoRecover()

					// Verify the contents of the request
					Expect(req.URL.EscapedPath()).To(Equal(getTrustStoreCertificatePath))
					Expect(req.Method).To(Equal("GET"))

					Expect(req.Header["Accept-Language"]).ToNot(BeNil())
					Expect(req.Header["Accept-Language"][0]).To(Equal(fmt.Sprintf("%v", "en-US,en;q=0.5")))
					// Set mock response
					res.Header().Set("Content-type", "application/json")
					res.WriteHeader(200)
					fmt.Fprintf(res, "%s", `{"id": "ID", "label": "Label", "certificate_type": "trust_store", "fingerprint_sha256": "FingerprintSha256", "subject_dn": "SubjectDn", "subject_cn": "SubjectCn", "issuer_dn": "IssuerDn", "issuer_cn": "IssuerCn", "issued": "2019-01-01T12:00:00.000Z", "expiry": "2019-01-01T12:00:00.000Z", "trusted": false, "href": "Href"}`)
				}))
			})
			It(`Invoke GetTrustStoreCertificate successfully`, func() {
				mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1(&mqcloudv1.MqcloudV1Options{
					URL:            testServer.URL,
					Authenticator:  &core.NoAuthAuthenticator{},
					AcceptLanguage: core.StringPtr(acceptLanguage),
				})
				Expect(serviceErr).To(BeNil())
				Expect(mqcloudService).ToNot(BeNil())

				// Invoke operation with nil options model (negative test)
				result, response, operationErr := mqcloudService.GetTrustStoreCertificate(nil)
				Expect(operationErr).NotTo(BeNil())
				Expect(response).To(BeNil())
				Expect(result).To(BeNil())

				// Construct an instance of the GetTrustStoreCertificateOptions model
				getTrustStoreCertificateOptionsModel := new(mqcloudv1.GetTrustStoreCertificateOptions)
				getTrustStoreCertificateOptionsModel.ServiceInstanceGuid = core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				getTrustStoreCertificateOptionsModel.QueueManagerID = core.StringPtr("b8e1aeda078009cf3db74e90d5d42328")
				getTrustStoreCertificateOptionsModel.CertificateID = core.StringPtr("9b7d1e723af8233")
				getTrustStoreCertificateOptionsModel.Headers = map[string]string{"x-custom-header": "x-custom-value"}

				// Invoke operation with valid options model (positive test)
				result, response, operationErr = mqcloudService.GetTrustStoreCertificate(getTrustStoreCertificateOptionsModel)
				Expect(operationErr).To(BeNil())
				Expect(response).ToNot(BeNil())
				Expect(result).ToNot(BeNil())

			})
			It(`Invoke GetTrustStoreCertificate with error: Operation validation and request error`, func() {
				mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1(&mqcloudv1.MqcloudV1Options{
					URL:            testServer.URL,
					Authenticator:  &core.NoAuthAuthenticator{},
					AcceptLanguage: core.StringPtr(acceptLanguage),
				})
				Expect(serviceErr).To(BeNil())
				Expect(mqcloudService).ToNot(BeNil())

				// Construct an instance of the GetTrustStoreCertificateOptions model
				getTrustStoreCertificateOptionsModel := new(mqcloudv1.GetTrustStoreCertificateOptions)
				getTrustStoreCertificateOptionsModel.ServiceInstanceGuid = core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				getTrustStoreCertificateOptionsModel.QueueManagerID = core.StringPtr("b8e1aeda078009cf3db74e90d5d42328")
				getTrustStoreCertificateOptionsModel.CertificateID = core.StringPtr("9b7d1e723af8233")
				getTrustStoreCertificateOptionsModel.Headers = map[string]string{"x-custom-header": "x-custom-value"}
				// Invoke operation with empty URL (negative test)
				err := mqcloudService.SetServiceURL("")
				Expect(err).To(BeNil())
				result, response, operationErr := mqcloudService.GetTrustStoreCertificate(getTrustStoreCertificateOptionsModel)
				Expect(operationErr).ToNot(BeNil())
				Expect(operationErr.Error()).To(ContainSubstring(core.ERRORMSG_SERVICE_URL_MISSING))
				Expect(response).To(BeNil())
				Expect(result).To(BeNil())
				// Construct a second instance of the GetTrustStoreCertificateOptions model with no property values
				getTrustStoreCertificateOptionsModelNew := new(mqcloudv1.GetTrustStoreCertificateOptions)
				// Invoke operation with invalid model (negative test)
				result, response, operationErr = mqcloudService.GetTrustStoreCertificate(getTrustStoreCertificateOptionsModelNew)
				Expect(operationErr).ToNot(BeNil())
				Expect(response).To(BeNil())
				Expect(result).To(BeNil())
			})
			AfterEach(func() {
				testServer.Close()
			})
		})
		Context(`Using mock server endpoint with missing response body`, func() {
			BeforeEach(func() {
				testServer = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
					defer GinkgoRecover()

					// Set success status code with no respoonse body
					res.WriteHeader(200)
				}))
			})
			It(`Invoke GetTrustStoreCertificate successfully`, func() {
				mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1(&mqcloudv1.MqcloudV1Options{
					URL:            testServer.URL,
					Authenticator:  &core.NoAuthAuthenticator{},
					AcceptLanguage: core.StringPtr(acceptLanguage),
				})
				Expect(serviceErr).To(BeNil())
				Expect(mqcloudService).ToNot(BeNil())

				// Construct an instance of the GetTrustStoreCertificateOptions model
				getTrustStoreCertificateOptionsModel := new(mqcloudv1.GetTrustStoreCertificateOptions)
				getTrustStoreCertificateOptionsModel.ServiceInstanceGuid = core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				getTrustStoreCertificateOptionsModel.QueueManagerID = core.StringPtr("b8e1aeda078009cf3db74e90d5d42328")
				getTrustStoreCertificateOptionsModel.CertificateID = core.StringPtr("9b7d1e723af8233")
				getTrustStoreCertificateOptionsModel.Headers = map[string]string{"x-custom-header": "x-custom-value"}

				// Invoke operation
				result, response, operationErr := mqcloudService.GetTrustStoreCertificate(getTrustStoreCertificateOptionsModel)
				Expect(operationErr).To(BeNil())
				Expect(response).ToNot(BeNil())

				// Verify a nil result
				Expect(result).To(BeNil())
			})
			AfterEach(func() {
				testServer.Close()
			})
		})
	})
	Describe(`DeleteTrustStoreCertificate(deleteTrustStoreCertificateOptions *DeleteTrustStoreCertificateOptions)`, func() {
		acceptLanguage := "en-US,en;q=0.5"
		deleteTrustStoreCertificatePath := "/v1/a2b4d4bc-dadb-4637-bcec-9b7d1e723af8/queue_managers/b8e1aeda078009cf3db74e90d5d42328/certificates/trust_store/9b7d1e723af8233"
		Context(`Using mock server endpoint`, func() {
			BeforeEach(func() {
				testServer = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
					defer GinkgoRecover()

					// Verify the contents of the request
					Expect(req.URL.EscapedPath()).To(Equal(deleteTrustStoreCertificatePath))
					Expect(req.Method).To(Equal("DELETE"))

					Expect(req.Header["Accept-Language"]).ToNot(BeNil())
					Expect(req.Header["Accept-Language"][0]).To(Equal(fmt.Sprintf("%v", "en-US,en;q=0.5")))
					res.WriteHeader(204)
				}))
			})
			It(`Invoke DeleteTrustStoreCertificate successfully`, func() {
				mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1(&mqcloudv1.MqcloudV1Options{
					URL:            testServer.URL,
					Authenticator:  &core.NoAuthAuthenticator{},
					AcceptLanguage: core.StringPtr(acceptLanguage),
				})
				Expect(serviceErr).To(BeNil())
				Expect(mqcloudService).ToNot(BeNil())

				// Invoke operation with nil options model (negative test)
				response, operationErr := mqcloudService.DeleteTrustStoreCertificate(nil)
				Expect(operationErr).NotTo(BeNil())
				Expect(response).To(BeNil())

				// Construct an instance of the DeleteTrustStoreCertificateOptions model
				deleteTrustStoreCertificateOptionsModel := new(mqcloudv1.DeleteTrustStoreCertificateOptions)
				deleteTrustStoreCertificateOptionsModel.ServiceInstanceGuid = core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				deleteTrustStoreCertificateOptionsModel.QueueManagerID = core.StringPtr("b8e1aeda078009cf3db74e90d5d42328")
				deleteTrustStoreCertificateOptionsModel.CertificateID = core.StringPtr("9b7d1e723af8233")
				deleteTrustStoreCertificateOptionsModel.Headers = map[string]string{"x-custom-header": "x-custom-value"}

				// Invoke operation with valid options model (positive test)
				response, operationErr = mqcloudService.DeleteTrustStoreCertificate(deleteTrustStoreCertificateOptionsModel)
				Expect(operationErr).To(BeNil())
				Expect(response).ToNot(BeNil())
			})
			It(`Invoke DeleteTrustStoreCertificate with error: Operation validation and request error`, func() {
				mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1(&mqcloudv1.MqcloudV1Options{
					URL:            testServer.URL,
					Authenticator:  &core.NoAuthAuthenticator{},
					AcceptLanguage: core.StringPtr(acceptLanguage),
				})
				Expect(serviceErr).To(BeNil())
				Expect(mqcloudService).ToNot(BeNil())

				// Construct an instance of the DeleteTrustStoreCertificateOptions model
				deleteTrustStoreCertificateOptionsModel := new(mqcloudv1.DeleteTrustStoreCertificateOptions)
				deleteTrustStoreCertificateOptionsModel.ServiceInstanceGuid = core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				deleteTrustStoreCertificateOptionsModel.QueueManagerID = core.StringPtr("b8e1aeda078009cf3db74e90d5d42328")
				deleteTrustStoreCertificateOptionsModel.CertificateID = core.StringPtr("9b7d1e723af8233")
				deleteTrustStoreCertificateOptionsModel.Headers = map[string]string{"x-custom-header": "x-custom-value"}
				// Invoke operation with empty URL (negative test)
				err := mqcloudService.SetServiceURL("")
				Expect(err).To(BeNil())
				response, operationErr := mqcloudService.DeleteTrustStoreCertificate(deleteTrustStoreCertificateOptionsModel)
				Expect(operationErr).ToNot(BeNil())
				Expect(operationErr.Error()).To(ContainSubstring(core.ERRORMSG_SERVICE_URL_MISSING))
				Expect(response).To(BeNil())
				// Construct a second instance of the DeleteTrustStoreCertificateOptions model with no property values
				deleteTrustStoreCertificateOptionsModelNew := new(mqcloudv1.DeleteTrustStoreCertificateOptions)
				// Invoke operation with invalid model (negative test)
				response, operationErr = mqcloudService.DeleteTrustStoreCertificate(deleteTrustStoreCertificateOptionsModelNew)
				Expect(operationErr).ToNot(BeNil())
				Expect(response).To(BeNil())
			})
			AfterEach(func() {
				testServer.Close()
			})
		})
	})
	Describe(`DownloadTrustStoreCertificate(downloadTrustStoreCertificateOptions *DownloadTrustStoreCertificateOptions)`, func() {
		acceptLanguage := "en-US,en;q=0.5"
		downloadTrustStoreCertificatePath := "/v1/a2b4d4bc-dadb-4637-bcec-9b7d1e723af8/queue_managers/b8e1aeda078009cf3db74e90d5d42328/certificates/trust_store/9b7d1e723af8233/download"
		Context(`Using mock server endpoint with timeout`, func() {
			BeforeEach(func() {
				testServer = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
					defer GinkgoRecover()

					// Verify the contents of the request
					Expect(req.URL.EscapedPath()).To(Equal(downloadTrustStoreCertificatePath))
					Expect(req.Method).To(Equal("GET"))

					Expect(req.Header["Accept-Language"]).ToNot(BeNil())
					Expect(req.Header["Accept-Language"][0]).To(Equal(fmt.Sprintf("%v", "en-US,en;q=0.5")))
					// Sleep a short time to support a timeout test
					time.Sleep(100 * time.Millisecond)

					// Set mock response
					res.Header().Set("Content-type", "application/octet-stream")
					res.WriteHeader(200)
					fmt.Fprintf(res, "%s", `This is a mock binary response.`)
				}))
			})
			It(`Invoke DownloadTrustStoreCertificate successfully with retries`, func() {
				mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1(&mqcloudv1.MqcloudV1Options{
					URL:            testServer.URL,
					Authenticator:  &core.NoAuthAuthenticator{},
					AcceptLanguage: core.StringPtr(acceptLanguage),
				})
				Expect(serviceErr).To(BeNil())
				Expect(mqcloudService).ToNot(BeNil())
				mqcloudService.EnableRetries(0, 0)

				// Construct an instance of the DownloadTrustStoreCertificateOptions model
				downloadTrustStoreCertificateOptionsModel := new(mqcloudv1.DownloadTrustStoreCertificateOptions)
				downloadTrustStoreCertificateOptionsModel.ServiceInstanceGuid = core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				downloadTrustStoreCertificateOptionsModel.QueueManagerID = core.StringPtr("b8e1aeda078009cf3db74e90d5d42328")
				downloadTrustStoreCertificateOptionsModel.CertificateID = core.StringPtr("9b7d1e723af8233")
				downloadTrustStoreCertificateOptionsModel.Headers = map[string]string{"x-custom-header": "x-custom-value"}

				// Invoke operation with a Context to test a timeout error
				ctx, cancelFunc := context.WithTimeout(context.Background(), 80*time.Millisecond)
				defer cancelFunc()
				_, _, operationErr := mqcloudService.DownloadTrustStoreCertificateWithContext(ctx, downloadTrustStoreCertificateOptionsModel)
				Expect(operationErr).ToNot(BeNil())
				Expect(operationErr.Error()).To(ContainSubstring("deadline exceeded"))

				// Disable retries and test again
				mqcloudService.DisableRetries()
				result, response, operationErr := mqcloudService.DownloadTrustStoreCertificate(downloadTrustStoreCertificateOptionsModel)
				Expect(operationErr).To(BeNil())
				Expect(response).ToNot(BeNil())
				Expect(result).ToNot(BeNil())

				// Re-test the timeout error with retries disabled
				ctx, cancelFunc2 := context.WithTimeout(context.Background(), 80*time.Millisecond)
				defer cancelFunc2()
				_, _, operationErr = mqcloudService.DownloadTrustStoreCertificateWithContext(ctx, downloadTrustStoreCertificateOptionsModel)
				Expect(operationErr).ToNot(BeNil())
				Expect(operationErr.Error()).To(ContainSubstring("deadline exceeded"))
			})
			AfterEach(func() {
				testServer.Close()
			})
		})
		Context(`Using mock server endpoint`, func() {
			BeforeEach(func() {
				testServer = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
					defer GinkgoRecover()

					// Verify the contents of the request
					Expect(req.URL.EscapedPath()).To(Equal(downloadTrustStoreCertificatePath))
					Expect(req.Method).To(Equal("GET"))

					Expect(req.Header["Accept-Language"]).ToNot(BeNil())
					Expect(req.Header["Accept-Language"][0]).To(Equal(fmt.Sprintf("%v", "en-US,en;q=0.5")))
					// Set mock response
					res.Header().Set("Content-type", "application/octet-stream")
					res.WriteHeader(200)
					fmt.Fprintf(res, "%s", `This is a mock binary response.`)
				}))
			})
			It(`Invoke DownloadTrustStoreCertificate successfully`, func() {
				mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1(&mqcloudv1.MqcloudV1Options{
					URL:            testServer.URL,
					Authenticator:  &core.NoAuthAuthenticator{},
					AcceptLanguage: core.StringPtr(acceptLanguage),
				})
				Expect(serviceErr).To(BeNil())
				Expect(mqcloudService).ToNot(BeNil())

				// Invoke operation with nil options model (negative test)
				result, response, operationErr := mqcloudService.DownloadTrustStoreCertificate(nil)
				Expect(operationErr).NotTo(BeNil())
				Expect(response).To(BeNil())
				Expect(result).To(BeNil())

				// Construct an instance of the DownloadTrustStoreCertificateOptions model
				downloadTrustStoreCertificateOptionsModel := new(mqcloudv1.DownloadTrustStoreCertificateOptions)
				downloadTrustStoreCertificateOptionsModel.ServiceInstanceGuid = core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				downloadTrustStoreCertificateOptionsModel.QueueManagerID = core.StringPtr("b8e1aeda078009cf3db74e90d5d42328")
				downloadTrustStoreCertificateOptionsModel.CertificateID = core.StringPtr("9b7d1e723af8233")
				downloadTrustStoreCertificateOptionsModel.Headers = map[string]string{"x-custom-header": "x-custom-value"}

				// Invoke operation with valid options model (positive test)
				result, response, operationErr = mqcloudService.DownloadTrustStoreCertificate(downloadTrustStoreCertificateOptionsModel)
				Expect(operationErr).To(BeNil())
				Expect(response).ToNot(BeNil())
				Expect(result).ToNot(BeNil())

			})
			It(`Invoke DownloadTrustStoreCertificate with error: Operation validation and request error`, func() {
				mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1(&mqcloudv1.MqcloudV1Options{
					URL:            testServer.URL,
					Authenticator:  &core.NoAuthAuthenticator{},
					AcceptLanguage: core.StringPtr(acceptLanguage),
				})
				Expect(serviceErr).To(BeNil())
				Expect(mqcloudService).ToNot(BeNil())

				// Construct an instance of the DownloadTrustStoreCertificateOptions model
				downloadTrustStoreCertificateOptionsModel := new(mqcloudv1.DownloadTrustStoreCertificateOptions)
				downloadTrustStoreCertificateOptionsModel.ServiceInstanceGuid = core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				downloadTrustStoreCertificateOptionsModel.QueueManagerID = core.StringPtr("b8e1aeda078009cf3db74e90d5d42328")
				downloadTrustStoreCertificateOptionsModel.CertificateID = core.StringPtr("9b7d1e723af8233")
				downloadTrustStoreCertificateOptionsModel.Headers = map[string]string{"x-custom-header": "x-custom-value"}
				// Invoke operation with empty URL (negative test)
				err := mqcloudService.SetServiceURL("")
				Expect(err).To(BeNil())
				result, response, operationErr := mqcloudService.DownloadTrustStoreCertificate(downloadTrustStoreCertificateOptionsModel)
				Expect(operationErr).ToNot(BeNil())
				Expect(operationErr.Error()).To(ContainSubstring(core.ERRORMSG_SERVICE_URL_MISSING))
				Expect(response).To(BeNil())
				Expect(result).To(BeNil())
				// Construct a second instance of the DownloadTrustStoreCertificateOptions model with no property values
				downloadTrustStoreCertificateOptionsModelNew := new(mqcloudv1.DownloadTrustStoreCertificateOptions)
				// Invoke operation with invalid model (negative test)
				result, response, operationErr = mqcloudService.DownloadTrustStoreCertificate(downloadTrustStoreCertificateOptionsModelNew)
				Expect(operationErr).ToNot(BeNil())
				Expect(response).To(BeNil())
				Expect(result).To(BeNil())
			})
			AfterEach(func() {
				testServer.Close()
			})
		})
		Context(`Using mock server endpoint with missing response body`, func() {
			BeforeEach(func() {
				testServer = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
					defer GinkgoRecover()

					// Set success status code with no respoonse body
					res.WriteHeader(200)
				}))
			})
			It(`Invoke DownloadTrustStoreCertificate successfully`, func() {
				mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1(&mqcloudv1.MqcloudV1Options{
					URL:            testServer.URL,
					Authenticator:  &core.NoAuthAuthenticator{},
					AcceptLanguage: core.StringPtr(acceptLanguage),
				})
				Expect(serviceErr).To(BeNil())
				Expect(mqcloudService).ToNot(BeNil())

				// Construct an instance of the DownloadTrustStoreCertificateOptions model
				downloadTrustStoreCertificateOptionsModel := new(mqcloudv1.DownloadTrustStoreCertificateOptions)
				downloadTrustStoreCertificateOptionsModel.ServiceInstanceGuid = core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				downloadTrustStoreCertificateOptionsModel.QueueManagerID = core.StringPtr("b8e1aeda078009cf3db74e90d5d42328")
				downloadTrustStoreCertificateOptionsModel.CertificateID = core.StringPtr("9b7d1e723af8233")
				downloadTrustStoreCertificateOptionsModel.Headers = map[string]string{"x-custom-header": "x-custom-value"}

				// Invoke operation
				result, response, operationErr := mqcloudService.DownloadTrustStoreCertificate(downloadTrustStoreCertificateOptionsModel)
				Expect(operationErr).To(BeNil())
				Expect(response).ToNot(BeNil())

				// Verify empty byte buffer.
				Expect(result).ToNot(BeNil())
				buffer, operationErr := io.ReadAll(result)
				Expect(operationErr).To(BeNil())
				Expect(buffer).ToNot(BeNil())
				Expect(len(buffer)).To(Equal(0))
			})
			AfterEach(func() {
				testServer.Close()
			})
		})
	})
	Describe(`CreateKeyStorePemCertificate(createKeyStorePemCertificateOptions *CreateKeyStorePemCertificateOptions) - Operation response error`, func() {
		acceptLanguage := "en-US,en;q=0.5"
		createKeyStorePemCertificatePath := "/v1/a2b4d4bc-dadb-4637-bcec-9b7d1e723af8/queue_managers/b8e1aeda078009cf3db74e90d5d42328/certificates/key_store"
		Context(`Using mock server endpoint with invalid JSON response`, func() {
			BeforeEach(func() {
				testServer = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
					defer GinkgoRecover()

					// Verify the contents of the request
					Expect(req.URL.EscapedPath()).To(Equal(createKeyStorePemCertificatePath))
					Expect(req.Method).To(Equal("POST"))
					Expect(req.Header["Accept-Language"]).ToNot(BeNil())
					Expect(req.Header["Accept-Language"][0]).To(Equal(fmt.Sprintf("%v", "en-US,en;q=0.5")))
					res.Header().Set("Content-type", "application/json")
					res.WriteHeader(201)
					fmt.Fprint(res, `} this is not valid json {`)
				}))
			})
			It(`Invoke CreateKeyStorePemCertificate with error: Operation response processing error`, func() {
				mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1(&mqcloudv1.MqcloudV1Options{
					URL:            testServer.URL,
					Authenticator:  &core.NoAuthAuthenticator{},
					AcceptLanguage: core.StringPtr(acceptLanguage),
				})
				Expect(serviceErr).To(BeNil())
				Expect(mqcloudService).ToNot(BeNil())

				// Construct an instance of the CreateKeyStorePemCertificateOptions model
				createKeyStorePemCertificateOptionsModel := new(mqcloudv1.CreateKeyStorePemCertificateOptions)
				createKeyStorePemCertificateOptionsModel.ServiceInstanceGuid = core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				createKeyStorePemCertificateOptionsModel.QueueManagerID = core.StringPtr("b8e1aeda078009cf3db74e90d5d42328")
				createKeyStorePemCertificateOptionsModel.Label = core.StringPtr("certlabel")
				createKeyStorePemCertificateOptionsModel.CertificateFile = CreateMockReader("This is a mock file.")
				createKeyStorePemCertificateOptionsModel.Headers = map[string]string{"x-custom-header": "x-custom-value"}
				// Expect response parsing to fail since we are receiving a text/plain response
				result, response, operationErr := mqcloudService.CreateKeyStorePemCertificate(createKeyStorePemCertificateOptionsModel)
				Expect(operationErr).ToNot(BeNil())
				Expect(response).ToNot(BeNil())
				Expect(result).To(BeNil())

				// Enable retries and test again
				mqcloudService.EnableRetries(0, 0)
				result, response, operationErr = mqcloudService.CreateKeyStorePemCertificate(createKeyStorePemCertificateOptionsModel)
				Expect(operationErr).ToNot(BeNil())
				Expect(response).ToNot(BeNil())
				Expect(result).To(BeNil())
			})
			AfterEach(func() {
				testServer.Close()
			})
		})
	})
	Describe(`CreateKeyStorePemCertificate(createKeyStorePemCertificateOptions *CreateKeyStorePemCertificateOptions)`, func() {
		acceptLanguage := "en-US,en;q=0.5"
		createKeyStorePemCertificatePath := "/v1/a2b4d4bc-dadb-4637-bcec-9b7d1e723af8/queue_managers/b8e1aeda078009cf3db74e90d5d42328/certificates/key_store"
		Context(`Using mock server endpoint with timeout`, func() {
			BeforeEach(func() {
				testServer = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
					defer GinkgoRecover()

					// Verify the contents of the request
					Expect(req.URL.EscapedPath()).To(Equal(createKeyStorePemCertificatePath))
					Expect(req.Method).To(Equal("POST"))

					Expect(req.Header["Accept-Language"]).ToNot(BeNil())
					Expect(req.Header["Accept-Language"][0]).To(Equal(fmt.Sprintf("%v", "en-US,en;q=0.5")))
					// Sleep a short time to support a timeout test
					time.Sleep(100 * time.Millisecond)

					// Set mock response
					res.Header().Set("Content-type", "application/json")
					res.WriteHeader(201)
					fmt.Fprintf(res, "%s", `{"id": "ID", "label": "Label", "certificate_type": "key_store", "fingerprint_sha256": "FingerprintSha256", "subject_dn": "SubjectDn", "subject_cn": "SubjectCn", "issuer_dn": "IssuerDn", "issuer_cn": "IssuerCn", "issued": "2019-01-01T12:00:00.000Z", "expiry": "2019-01-01T12:00:00.000Z", "is_default": false, "dns_names_total_count": 18, "dns_names": ["DnsNames"], "href": "Href", "config": {"ams": {"channels": [{"name": "Name"}]}}}`)
				}))
			})
			It(`Invoke CreateKeyStorePemCertificate successfully with retries`, func() {
				mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1(&mqcloudv1.MqcloudV1Options{
					URL:            testServer.URL,
					Authenticator:  &core.NoAuthAuthenticator{},
					AcceptLanguage: core.StringPtr(acceptLanguage),
				})
				Expect(serviceErr).To(BeNil())
				Expect(mqcloudService).ToNot(BeNil())
				mqcloudService.EnableRetries(0, 0)

				// Construct an instance of the CreateKeyStorePemCertificateOptions model
				createKeyStorePemCertificateOptionsModel := new(mqcloudv1.CreateKeyStorePemCertificateOptions)
				createKeyStorePemCertificateOptionsModel.ServiceInstanceGuid = core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				createKeyStorePemCertificateOptionsModel.QueueManagerID = core.StringPtr("b8e1aeda078009cf3db74e90d5d42328")
				createKeyStorePemCertificateOptionsModel.Label = core.StringPtr("certlabel")
				createKeyStorePemCertificateOptionsModel.CertificateFile = CreateMockReader("This is a mock file.")
				createKeyStorePemCertificateOptionsModel.Headers = map[string]string{"x-custom-header": "x-custom-value"}

				// Invoke operation with a Context to test a timeout error
				ctx, cancelFunc := context.WithTimeout(context.Background(), 80*time.Millisecond)
				defer cancelFunc()
				_, _, operationErr := mqcloudService.CreateKeyStorePemCertificateWithContext(ctx, createKeyStorePemCertificateOptionsModel)
				Expect(operationErr).ToNot(BeNil())
				Expect(operationErr.Error()).To(ContainSubstring("deadline exceeded"))

				// Disable retries and test again
				mqcloudService.DisableRetries()
				result, response, operationErr := mqcloudService.CreateKeyStorePemCertificate(createKeyStorePemCertificateOptionsModel)
				Expect(operationErr).To(BeNil())
				Expect(response).ToNot(BeNil())
				Expect(result).ToNot(BeNil())

				// Re-test the timeout error with retries disabled
				ctx, cancelFunc2 := context.WithTimeout(context.Background(), 80*time.Millisecond)
				defer cancelFunc2()
				_, _, operationErr = mqcloudService.CreateKeyStorePemCertificateWithContext(ctx, createKeyStorePemCertificateOptionsModel)
				Expect(operationErr).ToNot(BeNil())
				Expect(operationErr.Error()).To(ContainSubstring("deadline exceeded"))
			})
			AfterEach(func() {
				testServer.Close()
			})
		})
		Context(`Using mock server endpoint`, func() {
			BeforeEach(func() {
				testServer = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
					defer GinkgoRecover()

					// Verify the contents of the request
					Expect(req.URL.EscapedPath()).To(Equal(createKeyStorePemCertificatePath))
					Expect(req.Method).To(Equal("POST"))

					Expect(req.Header["Accept-Language"]).ToNot(BeNil())
					Expect(req.Header["Accept-Language"][0]).To(Equal(fmt.Sprintf("%v", "en-US,en;q=0.5")))
					// Set mock response
					res.Header().Set("Content-type", "application/json")
					res.WriteHeader(201)
					fmt.Fprintf(res, "%s", `{"id": "ID", "label": "Label", "certificate_type": "key_store", "fingerprint_sha256": "FingerprintSha256", "subject_dn": "SubjectDn", "subject_cn": "SubjectCn", "issuer_dn": "IssuerDn", "issuer_cn": "IssuerCn", "issued": "2019-01-01T12:00:00.000Z", "expiry": "2019-01-01T12:00:00.000Z", "is_default": false, "dns_names_total_count": 18, "dns_names": ["DnsNames"], "href": "Href", "config": {"ams": {"channels": [{"name": "Name"}]}}}`)
				}))
			})
			It(`Invoke CreateKeyStorePemCertificate successfully`, func() {
				mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1(&mqcloudv1.MqcloudV1Options{
					URL:            testServer.URL,
					Authenticator:  &core.NoAuthAuthenticator{},
					AcceptLanguage: core.StringPtr(acceptLanguage),
				})
				Expect(serviceErr).To(BeNil())
				Expect(mqcloudService).ToNot(BeNil())

				// Invoke operation with nil options model (negative test)
				result, response, operationErr := mqcloudService.CreateKeyStorePemCertificate(nil)
				Expect(operationErr).NotTo(BeNil())
				Expect(response).To(BeNil())
				Expect(result).To(BeNil())

				// Construct an instance of the CreateKeyStorePemCertificateOptions model
				createKeyStorePemCertificateOptionsModel := new(mqcloudv1.CreateKeyStorePemCertificateOptions)
				createKeyStorePemCertificateOptionsModel.ServiceInstanceGuid = core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				createKeyStorePemCertificateOptionsModel.QueueManagerID = core.StringPtr("b8e1aeda078009cf3db74e90d5d42328")
				createKeyStorePemCertificateOptionsModel.Label = core.StringPtr("certlabel")
				createKeyStorePemCertificateOptionsModel.CertificateFile = CreateMockReader("This is a mock file.")
				createKeyStorePemCertificateOptionsModel.Headers = map[string]string{"x-custom-header": "x-custom-value"}

				// Invoke operation with valid options model (positive test)
				result, response, operationErr = mqcloudService.CreateKeyStorePemCertificate(createKeyStorePemCertificateOptionsModel)
				Expect(operationErr).To(BeNil())
				Expect(response).ToNot(BeNil())
				Expect(result).ToNot(BeNil())

			})
			It(`Invoke CreateKeyStorePemCertificate with error: Operation validation and request error`, func() {
				mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1(&mqcloudv1.MqcloudV1Options{
					URL:            testServer.URL,
					Authenticator:  &core.NoAuthAuthenticator{},
					AcceptLanguage: core.StringPtr(acceptLanguage),
				})
				Expect(serviceErr).To(BeNil())
				Expect(mqcloudService).ToNot(BeNil())

				// Construct an instance of the CreateKeyStorePemCertificateOptions model
				createKeyStorePemCertificateOptionsModel := new(mqcloudv1.CreateKeyStorePemCertificateOptions)
				createKeyStorePemCertificateOptionsModel.ServiceInstanceGuid = core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				createKeyStorePemCertificateOptionsModel.QueueManagerID = core.StringPtr("b8e1aeda078009cf3db74e90d5d42328")
				createKeyStorePemCertificateOptionsModel.Label = core.StringPtr("certlabel")
				createKeyStorePemCertificateOptionsModel.CertificateFile = CreateMockReader("This is a mock file.")
				createKeyStorePemCertificateOptionsModel.Headers = map[string]string{"x-custom-header": "x-custom-value"}
				// Invoke operation with empty URL (negative test)
				err := mqcloudService.SetServiceURL("")
				Expect(err).To(BeNil())
				result, response, operationErr := mqcloudService.CreateKeyStorePemCertificate(createKeyStorePemCertificateOptionsModel)
				Expect(operationErr).ToNot(BeNil())
				Expect(operationErr.Error()).To(ContainSubstring(core.ERRORMSG_SERVICE_URL_MISSING))
				Expect(response).To(BeNil())
				Expect(result).To(BeNil())
				// Construct a second instance of the CreateKeyStorePemCertificateOptions model with no property values
				createKeyStorePemCertificateOptionsModelNew := new(mqcloudv1.CreateKeyStorePemCertificateOptions)
				// Invoke operation with invalid model (negative test)
				result, response, operationErr = mqcloudService.CreateKeyStorePemCertificate(createKeyStorePemCertificateOptionsModelNew)
				Expect(operationErr).ToNot(BeNil())
				Expect(response).To(BeNil())
				Expect(result).To(BeNil())
			})
			AfterEach(func() {
				testServer.Close()
			})
		})
		Context(`Using mock server endpoint with missing response body`, func() {
			BeforeEach(func() {
				testServer = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
					defer GinkgoRecover()

					// Set success status code with no respoonse body
					res.WriteHeader(201)
				}))
			})
			It(`Invoke CreateKeyStorePemCertificate successfully`, func() {
				mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1(&mqcloudv1.MqcloudV1Options{
					URL:            testServer.URL,
					Authenticator:  &core.NoAuthAuthenticator{},
					AcceptLanguage: core.StringPtr(acceptLanguage),
				})
				Expect(serviceErr).To(BeNil())
				Expect(mqcloudService).ToNot(BeNil())

				// Construct an instance of the CreateKeyStorePemCertificateOptions model
				createKeyStorePemCertificateOptionsModel := new(mqcloudv1.CreateKeyStorePemCertificateOptions)
				createKeyStorePemCertificateOptionsModel.ServiceInstanceGuid = core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				createKeyStorePemCertificateOptionsModel.QueueManagerID = core.StringPtr("b8e1aeda078009cf3db74e90d5d42328")
				createKeyStorePemCertificateOptionsModel.Label = core.StringPtr("certlabel")
				createKeyStorePemCertificateOptionsModel.CertificateFile = CreateMockReader("This is a mock file.")
				createKeyStorePemCertificateOptionsModel.Headers = map[string]string{"x-custom-header": "x-custom-value"}

				// Invoke operation
				result, response, operationErr := mqcloudService.CreateKeyStorePemCertificate(createKeyStorePemCertificateOptionsModel)
				Expect(operationErr).To(BeNil())
				Expect(response).ToNot(BeNil())

				// Verify a nil result
				Expect(result).To(BeNil())
			})
			AfterEach(func() {
				testServer.Close()
			})
		})
	})
	Describe(`ListKeyStoreCertificates(listKeyStoreCertificatesOptions *ListKeyStoreCertificatesOptions) - Operation response error`, func() {
		acceptLanguage := "en-US,en;q=0.5"
		listKeyStoreCertificatesPath := "/v1/a2b4d4bc-dadb-4637-bcec-9b7d1e723af8/queue_managers/b8e1aeda078009cf3db74e90d5d42328/certificates/key_store"
		Context(`Using mock server endpoint with invalid JSON response`, func() {
			BeforeEach(func() {
				testServer = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
					defer GinkgoRecover()

					// Verify the contents of the request
					Expect(req.URL.EscapedPath()).To(Equal(listKeyStoreCertificatesPath))
					Expect(req.Method).To(Equal("GET"))
					Expect(req.Header["Accept-Language"]).ToNot(BeNil())
					Expect(req.Header["Accept-Language"][0]).To(Equal(fmt.Sprintf("%v", "en-US,en;q=0.5")))
					res.Header().Set("Content-type", "application/json")
					res.WriteHeader(200)
					fmt.Fprint(res, `} this is not valid json {`)
				}))
			})
			It(`Invoke ListKeyStoreCertificates with error: Operation response processing error`, func() {
				mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1(&mqcloudv1.MqcloudV1Options{
					URL:            testServer.URL,
					Authenticator:  &core.NoAuthAuthenticator{},
					AcceptLanguage: core.StringPtr(acceptLanguage),
				})
				Expect(serviceErr).To(BeNil())
				Expect(mqcloudService).ToNot(BeNil())

				// Construct an instance of the ListKeyStoreCertificatesOptions model
				listKeyStoreCertificatesOptionsModel := new(mqcloudv1.ListKeyStoreCertificatesOptions)
				listKeyStoreCertificatesOptionsModel.ServiceInstanceGuid = core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				listKeyStoreCertificatesOptionsModel.QueueManagerID = core.StringPtr("b8e1aeda078009cf3db74e90d5d42328")
				listKeyStoreCertificatesOptionsModel.Headers = map[string]string{"x-custom-header": "x-custom-value"}
				// Expect response parsing to fail since we are receiving a text/plain response
				result, response, operationErr := mqcloudService.ListKeyStoreCertificates(listKeyStoreCertificatesOptionsModel)
				Expect(operationErr).ToNot(BeNil())
				Expect(response).ToNot(BeNil())
				Expect(result).To(BeNil())

				// Enable retries and test again
				mqcloudService.EnableRetries(0, 0)
				result, response, operationErr = mqcloudService.ListKeyStoreCertificates(listKeyStoreCertificatesOptionsModel)
				Expect(operationErr).ToNot(BeNil())
				Expect(response).ToNot(BeNil())
				Expect(result).To(BeNil())
			})
			AfterEach(func() {
				testServer.Close()
			})
		})
	})
	Describe(`ListKeyStoreCertificates(listKeyStoreCertificatesOptions *ListKeyStoreCertificatesOptions)`, func() {
		acceptLanguage := "en-US,en;q=0.5"
		listKeyStoreCertificatesPath := "/v1/a2b4d4bc-dadb-4637-bcec-9b7d1e723af8/queue_managers/b8e1aeda078009cf3db74e90d5d42328/certificates/key_store"
		Context(`Using mock server endpoint with timeout`, func() {
			BeforeEach(func() {
				testServer = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
					defer GinkgoRecover()

					// Verify the contents of the request
					Expect(req.URL.EscapedPath()).To(Equal(listKeyStoreCertificatesPath))
					Expect(req.Method).To(Equal("GET"))

					Expect(req.Header["Accept-Language"]).ToNot(BeNil())
					Expect(req.Header["Accept-Language"][0]).To(Equal(fmt.Sprintf("%v", "en-US,en;q=0.5")))
					// Sleep a short time to support a timeout test
					time.Sleep(100 * time.Millisecond)

					// Set mock response
					res.Header().Set("Content-type", "application/json")
					res.WriteHeader(200)
					fmt.Fprintf(res, "%s", `{"total_count": 1, "key_store": [{"id": "ID", "label": "Label", "certificate_type": "key_store", "fingerprint_sha256": "FingerprintSha256", "subject_dn": "SubjectDn", "subject_cn": "SubjectCn", "issuer_dn": "IssuerDn", "issuer_cn": "IssuerCn", "issued": "2019-01-01T12:00:00.000Z", "expiry": "2019-01-01T12:00:00.000Z", "is_default": false, "dns_names_total_count": 18, "dns_names": ["DnsNames"], "href": "Href", "config": {"ams": {"channels": [{"name": "Name"}]}}}]}`)
				}))
			})
			It(`Invoke ListKeyStoreCertificates successfully with retries`, func() {
				mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1(&mqcloudv1.MqcloudV1Options{
					URL:            testServer.URL,
					Authenticator:  &core.NoAuthAuthenticator{},
					AcceptLanguage: core.StringPtr(acceptLanguage),
				})
				Expect(serviceErr).To(BeNil())
				Expect(mqcloudService).ToNot(BeNil())
				mqcloudService.EnableRetries(0, 0)

				// Construct an instance of the ListKeyStoreCertificatesOptions model
				listKeyStoreCertificatesOptionsModel := new(mqcloudv1.ListKeyStoreCertificatesOptions)
				listKeyStoreCertificatesOptionsModel.ServiceInstanceGuid = core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				listKeyStoreCertificatesOptionsModel.QueueManagerID = core.StringPtr("b8e1aeda078009cf3db74e90d5d42328")
				listKeyStoreCertificatesOptionsModel.Headers = map[string]string{"x-custom-header": "x-custom-value"}

				// Invoke operation with a Context to test a timeout error
				ctx, cancelFunc := context.WithTimeout(context.Background(), 80*time.Millisecond)
				defer cancelFunc()
				_, _, operationErr := mqcloudService.ListKeyStoreCertificatesWithContext(ctx, listKeyStoreCertificatesOptionsModel)
				Expect(operationErr).ToNot(BeNil())
				Expect(operationErr.Error()).To(ContainSubstring("deadline exceeded"))

				// Disable retries and test again
				mqcloudService.DisableRetries()
				result, response, operationErr := mqcloudService.ListKeyStoreCertificates(listKeyStoreCertificatesOptionsModel)
				Expect(operationErr).To(BeNil())
				Expect(response).ToNot(BeNil())
				Expect(result).ToNot(BeNil())

				// Re-test the timeout error with retries disabled
				ctx, cancelFunc2 := context.WithTimeout(context.Background(), 80*time.Millisecond)
				defer cancelFunc2()
				_, _, operationErr = mqcloudService.ListKeyStoreCertificatesWithContext(ctx, listKeyStoreCertificatesOptionsModel)
				Expect(operationErr).ToNot(BeNil())
				Expect(operationErr.Error()).To(ContainSubstring("deadline exceeded"))
			})
			AfterEach(func() {
				testServer.Close()
			})
		})
		Context(`Using mock server endpoint`, func() {
			BeforeEach(func() {
				testServer = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
					defer GinkgoRecover()

					// Verify the contents of the request
					Expect(req.URL.EscapedPath()).To(Equal(listKeyStoreCertificatesPath))
					Expect(req.Method).To(Equal("GET"))

					Expect(req.Header["Accept-Language"]).ToNot(BeNil())
					Expect(req.Header["Accept-Language"][0]).To(Equal(fmt.Sprintf("%v", "en-US,en;q=0.5")))
					// Set mock response
					res.Header().Set("Content-type", "application/json")
					res.WriteHeader(200)
					fmt.Fprintf(res, "%s", `{"total_count": 1, "key_store": [{"id": "ID", "label": "Label", "certificate_type": "key_store", "fingerprint_sha256": "FingerprintSha256", "subject_dn": "SubjectDn", "subject_cn": "SubjectCn", "issuer_dn": "IssuerDn", "issuer_cn": "IssuerCn", "issued": "2019-01-01T12:00:00.000Z", "expiry": "2019-01-01T12:00:00.000Z", "is_default": false, "dns_names_total_count": 18, "dns_names": ["DnsNames"], "href": "Href", "config": {"ams": {"channels": [{"name": "Name"}]}}}]}`)
				}))
			})
			It(`Invoke ListKeyStoreCertificates successfully`, func() {
				mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1(&mqcloudv1.MqcloudV1Options{
					URL:            testServer.URL,
					Authenticator:  &core.NoAuthAuthenticator{},
					AcceptLanguage: core.StringPtr(acceptLanguage),
				})
				Expect(serviceErr).To(BeNil())
				Expect(mqcloudService).ToNot(BeNil())

				// Invoke operation with nil options model (negative test)
				result, response, operationErr := mqcloudService.ListKeyStoreCertificates(nil)
				Expect(operationErr).NotTo(BeNil())
				Expect(response).To(BeNil())
				Expect(result).To(BeNil())

				// Construct an instance of the ListKeyStoreCertificatesOptions model
				listKeyStoreCertificatesOptionsModel := new(mqcloudv1.ListKeyStoreCertificatesOptions)
				listKeyStoreCertificatesOptionsModel.ServiceInstanceGuid = core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				listKeyStoreCertificatesOptionsModel.QueueManagerID = core.StringPtr("b8e1aeda078009cf3db74e90d5d42328")
				listKeyStoreCertificatesOptionsModel.Headers = map[string]string{"x-custom-header": "x-custom-value"}

				// Invoke operation with valid options model (positive test)
				result, response, operationErr = mqcloudService.ListKeyStoreCertificates(listKeyStoreCertificatesOptionsModel)
				Expect(operationErr).To(BeNil())
				Expect(response).ToNot(BeNil())
				Expect(result).ToNot(BeNil())

			})
			It(`Invoke ListKeyStoreCertificates with error: Operation validation and request error`, func() {
				mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1(&mqcloudv1.MqcloudV1Options{
					URL:            testServer.URL,
					Authenticator:  &core.NoAuthAuthenticator{},
					AcceptLanguage: core.StringPtr(acceptLanguage),
				})
				Expect(serviceErr).To(BeNil())
				Expect(mqcloudService).ToNot(BeNil())

				// Construct an instance of the ListKeyStoreCertificatesOptions model
				listKeyStoreCertificatesOptionsModel := new(mqcloudv1.ListKeyStoreCertificatesOptions)
				listKeyStoreCertificatesOptionsModel.ServiceInstanceGuid = core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				listKeyStoreCertificatesOptionsModel.QueueManagerID = core.StringPtr("b8e1aeda078009cf3db74e90d5d42328")
				listKeyStoreCertificatesOptionsModel.Headers = map[string]string{"x-custom-header": "x-custom-value"}
				// Invoke operation with empty URL (negative test)
				err := mqcloudService.SetServiceURL("")
				Expect(err).To(BeNil())
				result, response, operationErr := mqcloudService.ListKeyStoreCertificates(listKeyStoreCertificatesOptionsModel)
				Expect(operationErr).ToNot(BeNil())
				Expect(operationErr.Error()).To(ContainSubstring(core.ERRORMSG_SERVICE_URL_MISSING))
				Expect(response).To(BeNil())
				Expect(result).To(BeNil())
				// Construct a second instance of the ListKeyStoreCertificatesOptions model with no property values
				listKeyStoreCertificatesOptionsModelNew := new(mqcloudv1.ListKeyStoreCertificatesOptions)
				// Invoke operation with invalid model (negative test)
				result, response, operationErr = mqcloudService.ListKeyStoreCertificates(listKeyStoreCertificatesOptionsModelNew)
				Expect(operationErr).ToNot(BeNil())
				Expect(response).To(BeNil())
				Expect(result).To(BeNil())
			})
			AfterEach(func() {
				testServer.Close()
			})
		})
		Context(`Using mock server endpoint with missing response body`, func() {
			BeforeEach(func() {
				testServer = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
					defer GinkgoRecover()

					// Set success status code with no respoonse body
					res.WriteHeader(200)
				}))
			})
			It(`Invoke ListKeyStoreCertificates successfully`, func() {
				mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1(&mqcloudv1.MqcloudV1Options{
					URL:            testServer.URL,
					Authenticator:  &core.NoAuthAuthenticator{},
					AcceptLanguage: core.StringPtr(acceptLanguage),
				})
				Expect(serviceErr).To(BeNil())
				Expect(mqcloudService).ToNot(BeNil())

				// Construct an instance of the ListKeyStoreCertificatesOptions model
				listKeyStoreCertificatesOptionsModel := new(mqcloudv1.ListKeyStoreCertificatesOptions)
				listKeyStoreCertificatesOptionsModel.ServiceInstanceGuid = core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				listKeyStoreCertificatesOptionsModel.QueueManagerID = core.StringPtr("b8e1aeda078009cf3db74e90d5d42328")
				listKeyStoreCertificatesOptionsModel.Headers = map[string]string{"x-custom-header": "x-custom-value"}

				// Invoke operation
				result, response, operationErr := mqcloudService.ListKeyStoreCertificates(listKeyStoreCertificatesOptionsModel)
				Expect(operationErr).To(BeNil())
				Expect(response).ToNot(BeNil())

				// Verify a nil result
				Expect(result).To(BeNil())
			})
			AfterEach(func() {
				testServer.Close()
			})
		})
	})
	Describe(`GetKeyStoreCertificate(getKeyStoreCertificateOptions *GetKeyStoreCertificateOptions) - Operation response error`, func() {
		acceptLanguage := "en-US,en;q=0.5"
		getKeyStoreCertificatePath := "/v1/a2b4d4bc-dadb-4637-bcec-9b7d1e723af8/queue_managers/b8e1aeda078009cf3db74e90d5d42328/certificates/key_store/9b7d1e723af8233"
		Context(`Using mock server endpoint with invalid JSON response`, func() {
			BeforeEach(func() {
				testServer = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
					defer GinkgoRecover()

					// Verify the contents of the request
					Expect(req.URL.EscapedPath()).To(Equal(getKeyStoreCertificatePath))
					Expect(req.Method).To(Equal("GET"))
					Expect(req.Header["Accept-Language"]).ToNot(BeNil())
					Expect(req.Header["Accept-Language"][0]).To(Equal(fmt.Sprintf("%v", "en-US,en;q=0.5")))
					res.Header().Set("Content-type", "application/json")
					res.WriteHeader(200)
					fmt.Fprint(res, `} this is not valid json {`)
				}))
			})
			It(`Invoke GetKeyStoreCertificate with error: Operation response processing error`, func() {
				mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1(&mqcloudv1.MqcloudV1Options{
					URL:            testServer.URL,
					Authenticator:  &core.NoAuthAuthenticator{},
					AcceptLanguage: core.StringPtr(acceptLanguage),
				})
				Expect(serviceErr).To(BeNil())
				Expect(mqcloudService).ToNot(BeNil())

				// Construct an instance of the GetKeyStoreCertificateOptions model
				getKeyStoreCertificateOptionsModel := new(mqcloudv1.GetKeyStoreCertificateOptions)
				getKeyStoreCertificateOptionsModel.ServiceInstanceGuid = core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				getKeyStoreCertificateOptionsModel.QueueManagerID = core.StringPtr("b8e1aeda078009cf3db74e90d5d42328")
				getKeyStoreCertificateOptionsModel.CertificateID = core.StringPtr("9b7d1e723af8233")
				getKeyStoreCertificateOptionsModel.Headers = map[string]string{"x-custom-header": "x-custom-value"}
				// Expect response parsing to fail since we are receiving a text/plain response
				result, response, operationErr := mqcloudService.GetKeyStoreCertificate(getKeyStoreCertificateOptionsModel)
				Expect(operationErr).ToNot(BeNil())
				Expect(response).ToNot(BeNil())
				Expect(result).To(BeNil())

				// Enable retries and test again
				mqcloudService.EnableRetries(0, 0)
				result, response, operationErr = mqcloudService.GetKeyStoreCertificate(getKeyStoreCertificateOptionsModel)
				Expect(operationErr).ToNot(BeNil())
				Expect(response).ToNot(BeNil())
				Expect(result).To(BeNil())
			})
			AfterEach(func() {
				testServer.Close()
			})
		})
	})
	Describe(`GetKeyStoreCertificate(getKeyStoreCertificateOptions *GetKeyStoreCertificateOptions)`, func() {
		acceptLanguage := "en-US,en;q=0.5"
		getKeyStoreCertificatePath := "/v1/a2b4d4bc-dadb-4637-bcec-9b7d1e723af8/queue_managers/b8e1aeda078009cf3db74e90d5d42328/certificates/key_store/9b7d1e723af8233"
		Context(`Using mock server endpoint with timeout`, func() {
			BeforeEach(func() {
				testServer = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
					defer GinkgoRecover()

					// Verify the contents of the request
					Expect(req.URL.EscapedPath()).To(Equal(getKeyStoreCertificatePath))
					Expect(req.Method).To(Equal("GET"))

					Expect(req.Header["Accept-Language"]).ToNot(BeNil())
					Expect(req.Header["Accept-Language"][0]).To(Equal(fmt.Sprintf("%v", "en-US,en;q=0.5")))
					// Sleep a short time to support a timeout test
					time.Sleep(100 * time.Millisecond)

					// Set mock response
					res.Header().Set("Content-type", "application/json")
					res.WriteHeader(200)
					fmt.Fprintf(res, "%s", `{"id": "ID", "label": "Label", "certificate_type": "key_store", "fingerprint_sha256": "FingerprintSha256", "subject_dn": "SubjectDn", "subject_cn": "SubjectCn", "issuer_dn": "IssuerDn", "issuer_cn": "IssuerCn", "issued": "2019-01-01T12:00:00.000Z", "expiry": "2019-01-01T12:00:00.000Z", "is_default": false, "dns_names_total_count": 18, "dns_names": ["DnsNames"], "href": "Href", "config": {"ams": {"channels": [{"name": "Name"}]}}}`)
				}))
			})
			It(`Invoke GetKeyStoreCertificate successfully with retries`, func() {
				mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1(&mqcloudv1.MqcloudV1Options{
					URL:            testServer.URL,
					Authenticator:  &core.NoAuthAuthenticator{},
					AcceptLanguage: core.StringPtr(acceptLanguage),
				})
				Expect(serviceErr).To(BeNil())
				Expect(mqcloudService).ToNot(BeNil())
				mqcloudService.EnableRetries(0, 0)

				// Construct an instance of the GetKeyStoreCertificateOptions model
				getKeyStoreCertificateOptionsModel := new(mqcloudv1.GetKeyStoreCertificateOptions)
				getKeyStoreCertificateOptionsModel.ServiceInstanceGuid = core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				getKeyStoreCertificateOptionsModel.QueueManagerID = core.StringPtr("b8e1aeda078009cf3db74e90d5d42328")
				getKeyStoreCertificateOptionsModel.CertificateID = core.StringPtr("9b7d1e723af8233")
				getKeyStoreCertificateOptionsModel.Headers = map[string]string{"x-custom-header": "x-custom-value"}

				// Invoke operation with a Context to test a timeout error
				ctx, cancelFunc := context.WithTimeout(context.Background(), 80*time.Millisecond)
				defer cancelFunc()
				_, _, operationErr := mqcloudService.GetKeyStoreCertificateWithContext(ctx, getKeyStoreCertificateOptionsModel)
				Expect(operationErr).ToNot(BeNil())
				Expect(operationErr.Error()).To(ContainSubstring("deadline exceeded"))

				// Disable retries and test again
				mqcloudService.DisableRetries()
				result, response, operationErr := mqcloudService.GetKeyStoreCertificate(getKeyStoreCertificateOptionsModel)
				Expect(operationErr).To(BeNil())
				Expect(response).ToNot(BeNil())
				Expect(result).ToNot(BeNil())

				// Re-test the timeout error with retries disabled
				ctx, cancelFunc2 := context.WithTimeout(context.Background(), 80*time.Millisecond)
				defer cancelFunc2()
				_, _, operationErr = mqcloudService.GetKeyStoreCertificateWithContext(ctx, getKeyStoreCertificateOptionsModel)
				Expect(operationErr).ToNot(BeNil())
				Expect(operationErr.Error()).To(ContainSubstring("deadline exceeded"))
			})
			AfterEach(func() {
				testServer.Close()
			})
		})
		Context(`Using mock server endpoint`, func() {
			BeforeEach(func() {
				testServer = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
					defer GinkgoRecover()

					// Verify the contents of the request
					Expect(req.URL.EscapedPath()).To(Equal(getKeyStoreCertificatePath))
					Expect(req.Method).To(Equal("GET"))

					Expect(req.Header["Accept-Language"]).ToNot(BeNil())
					Expect(req.Header["Accept-Language"][0]).To(Equal(fmt.Sprintf("%v", "en-US,en;q=0.5")))
					// Set mock response
					res.Header().Set("Content-type", "application/json")
					res.WriteHeader(200)
					fmt.Fprintf(res, "%s", `{"id": "ID", "label": "Label", "certificate_type": "key_store", "fingerprint_sha256": "FingerprintSha256", "subject_dn": "SubjectDn", "subject_cn": "SubjectCn", "issuer_dn": "IssuerDn", "issuer_cn": "IssuerCn", "issued": "2019-01-01T12:00:00.000Z", "expiry": "2019-01-01T12:00:00.000Z", "is_default": false, "dns_names_total_count": 18, "dns_names": ["DnsNames"], "href": "Href", "config": {"ams": {"channels": [{"name": "Name"}]}}}`)
				}))
			})
			It(`Invoke GetKeyStoreCertificate successfully`, func() {
				mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1(&mqcloudv1.MqcloudV1Options{
					URL:            testServer.URL,
					Authenticator:  &core.NoAuthAuthenticator{},
					AcceptLanguage: core.StringPtr(acceptLanguage),
				})
				Expect(serviceErr).To(BeNil())
				Expect(mqcloudService).ToNot(BeNil())

				// Invoke operation with nil options model (negative test)
				result, response, operationErr := mqcloudService.GetKeyStoreCertificate(nil)
				Expect(operationErr).NotTo(BeNil())
				Expect(response).To(BeNil())
				Expect(result).To(BeNil())

				// Construct an instance of the GetKeyStoreCertificateOptions model
				getKeyStoreCertificateOptionsModel := new(mqcloudv1.GetKeyStoreCertificateOptions)
				getKeyStoreCertificateOptionsModel.ServiceInstanceGuid = core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				getKeyStoreCertificateOptionsModel.QueueManagerID = core.StringPtr("b8e1aeda078009cf3db74e90d5d42328")
				getKeyStoreCertificateOptionsModel.CertificateID = core.StringPtr("9b7d1e723af8233")
				getKeyStoreCertificateOptionsModel.Headers = map[string]string{"x-custom-header": "x-custom-value"}

				// Invoke operation with valid options model (positive test)
				result, response, operationErr = mqcloudService.GetKeyStoreCertificate(getKeyStoreCertificateOptionsModel)
				Expect(operationErr).To(BeNil())
				Expect(response).ToNot(BeNil())
				Expect(result).ToNot(BeNil())

			})
			It(`Invoke GetKeyStoreCertificate with error: Operation validation and request error`, func() {
				mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1(&mqcloudv1.MqcloudV1Options{
					URL:            testServer.URL,
					Authenticator:  &core.NoAuthAuthenticator{},
					AcceptLanguage: core.StringPtr(acceptLanguage),
				})
				Expect(serviceErr).To(BeNil())
				Expect(mqcloudService).ToNot(BeNil())

				// Construct an instance of the GetKeyStoreCertificateOptions model
				getKeyStoreCertificateOptionsModel := new(mqcloudv1.GetKeyStoreCertificateOptions)
				getKeyStoreCertificateOptionsModel.ServiceInstanceGuid = core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				getKeyStoreCertificateOptionsModel.QueueManagerID = core.StringPtr("b8e1aeda078009cf3db74e90d5d42328")
				getKeyStoreCertificateOptionsModel.CertificateID = core.StringPtr("9b7d1e723af8233")
				getKeyStoreCertificateOptionsModel.Headers = map[string]string{"x-custom-header": "x-custom-value"}
				// Invoke operation with empty URL (negative test)
				err := mqcloudService.SetServiceURL("")
				Expect(err).To(BeNil())
				result, response, operationErr := mqcloudService.GetKeyStoreCertificate(getKeyStoreCertificateOptionsModel)
				Expect(operationErr).ToNot(BeNil())
				Expect(operationErr.Error()).To(ContainSubstring(core.ERRORMSG_SERVICE_URL_MISSING))
				Expect(response).To(BeNil())
				Expect(result).To(BeNil())
				// Construct a second instance of the GetKeyStoreCertificateOptions model with no property values
				getKeyStoreCertificateOptionsModelNew := new(mqcloudv1.GetKeyStoreCertificateOptions)
				// Invoke operation with invalid model (negative test)
				result, response, operationErr = mqcloudService.GetKeyStoreCertificate(getKeyStoreCertificateOptionsModelNew)
				Expect(operationErr).ToNot(BeNil())
				Expect(response).To(BeNil())
				Expect(result).To(BeNil())
			})
			AfterEach(func() {
				testServer.Close()
			})
		})
		Context(`Using mock server endpoint with missing response body`, func() {
			BeforeEach(func() {
				testServer = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
					defer GinkgoRecover()

					// Set success status code with no respoonse body
					res.WriteHeader(200)
				}))
			})
			It(`Invoke GetKeyStoreCertificate successfully`, func() {
				mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1(&mqcloudv1.MqcloudV1Options{
					URL:            testServer.URL,
					Authenticator:  &core.NoAuthAuthenticator{},
					AcceptLanguage: core.StringPtr(acceptLanguage),
				})
				Expect(serviceErr).To(BeNil())
				Expect(mqcloudService).ToNot(BeNil())

				// Construct an instance of the GetKeyStoreCertificateOptions model
				getKeyStoreCertificateOptionsModel := new(mqcloudv1.GetKeyStoreCertificateOptions)
				getKeyStoreCertificateOptionsModel.ServiceInstanceGuid = core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				getKeyStoreCertificateOptionsModel.QueueManagerID = core.StringPtr("b8e1aeda078009cf3db74e90d5d42328")
				getKeyStoreCertificateOptionsModel.CertificateID = core.StringPtr("9b7d1e723af8233")
				getKeyStoreCertificateOptionsModel.Headers = map[string]string{"x-custom-header": "x-custom-value"}

				// Invoke operation
				result, response, operationErr := mqcloudService.GetKeyStoreCertificate(getKeyStoreCertificateOptionsModel)
				Expect(operationErr).To(BeNil())
				Expect(response).ToNot(BeNil())

				// Verify a nil result
				Expect(result).To(BeNil())
			})
			AfterEach(func() {
				testServer.Close()
			})
		})
	})
	Describe(`DeleteKeyStoreCertificate(deleteKeyStoreCertificateOptions *DeleteKeyStoreCertificateOptions)`, func() {
		acceptLanguage := "en-US,en;q=0.5"
		deleteKeyStoreCertificatePath := "/v1/a2b4d4bc-dadb-4637-bcec-9b7d1e723af8/queue_managers/b8e1aeda078009cf3db74e90d5d42328/certificates/key_store/9b7d1e723af8233"
		Context(`Using mock server endpoint`, func() {
			BeforeEach(func() {
				testServer = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
					defer GinkgoRecover()

					// Verify the contents of the request
					Expect(req.URL.EscapedPath()).To(Equal(deleteKeyStoreCertificatePath))
					Expect(req.Method).To(Equal("DELETE"))

					Expect(req.Header["Accept-Language"]).ToNot(BeNil())
					Expect(req.Header["Accept-Language"][0]).To(Equal(fmt.Sprintf("%v", "en-US,en;q=0.5")))
					res.WriteHeader(204)
				}))
			})
			It(`Invoke DeleteKeyStoreCertificate successfully`, func() {
				mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1(&mqcloudv1.MqcloudV1Options{
					URL:            testServer.URL,
					Authenticator:  &core.NoAuthAuthenticator{},
					AcceptLanguage: core.StringPtr(acceptLanguage),
				})
				Expect(serviceErr).To(BeNil())
				Expect(mqcloudService).ToNot(BeNil())

				// Invoke operation with nil options model (negative test)
				response, operationErr := mqcloudService.DeleteKeyStoreCertificate(nil)
				Expect(operationErr).NotTo(BeNil())
				Expect(response).To(BeNil())

				// Construct an instance of the DeleteKeyStoreCertificateOptions model
				deleteKeyStoreCertificateOptionsModel := new(mqcloudv1.DeleteKeyStoreCertificateOptions)
				deleteKeyStoreCertificateOptionsModel.ServiceInstanceGuid = core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				deleteKeyStoreCertificateOptionsModel.QueueManagerID = core.StringPtr("b8e1aeda078009cf3db74e90d5d42328")
				deleteKeyStoreCertificateOptionsModel.CertificateID = core.StringPtr("9b7d1e723af8233")
				deleteKeyStoreCertificateOptionsModel.Headers = map[string]string{"x-custom-header": "x-custom-value"}

				// Invoke operation with valid options model (positive test)
				response, operationErr = mqcloudService.DeleteKeyStoreCertificate(deleteKeyStoreCertificateOptionsModel)
				Expect(operationErr).To(BeNil())
				Expect(response).ToNot(BeNil())
			})
			It(`Invoke DeleteKeyStoreCertificate with error: Operation validation and request error`, func() {
				mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1(&mqcloudv1.MqcloudV1Options{
					URL:            testServer.URL,
					Authenticator:  &core.NoAuthAuthenticator{},
					AcceptLanguage: core.StringPtr(acceptLanguage),
				})
				Expect(serviceErr).To(BeNil())
				Expect(mqcloudService).ToNot(BeNil())

				// Construct an instance of the DeleteKeyStoreCertificateOptions model
				deleteKeyStoreCertificateOptionsModel := new(mqcloudv1.DeleteKeyStoreCertificateOptions)
				deleteKeyStoreCertificateOptionsModel.ServiceInstanceGuid = core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				deleteKeyStoreCertificateOptionsModel.QueueManagerID = core.StringPtr("b8e1aeda078009cf3db74e90d5d42328")
				deleteKeyStoreCertificateOptionsModel.CertificateID = core.StringPtr("9b7d1e723af8233")
				deleteKeyStoreCertificateOptionsModel.Headers = map[string]string{"x-custom-header": "x-custom-value"}
				// Invoke operation with empty URL (negative test)
				err := mqcloudService.SetServiceURL("")
				Expect(err).To(BeNil())
				response, operationErr := mqcloudService.DeleteKeyStoreCertificate(deleteKeyStoreCertificateOptionsModel)
				Expect(operationErr).ToNot(BeNil())
				Expect(operationErr.Error()).To(ContainSubstring(core.ERRORMSG_SERVICE_URL_MISSING))
				Expect(response).To(BeNil())
				// Construct a second instance of the DeleteKeyStoreCertificateOptions model with no property values
				deleteKeyStoreCertificateOptionsModelNew := new(mqcloudv1.DeleteKeyStoreCertificateOptions)
				// Invoke operation with invalid model (negative test)
				response, operationErr = mqcloudService.DeleteKeyStoreCertificate(deleteKeyStoreCertificateOptionsModelNew)
				Expect(operationErr).ToNot(BeNil())
				Expect(response).To(BeNil())
			})
			AfterEach(func() {
				testServer.Close()
			})
		})
	})
	Describe(`DownloadKeyStoreCertificate(downloadKeyStoreCertificateOptions *DownloadKeyStoreCertificateOptions)`, func() {
		acceptLanguage := "en-US,en;q=0.5"
		downloadKeyStoreCertificatePath := "/v1/a2b4d4bc-dadb-4637-bcec-9b7d1e723af8/queue_managers/b8e1aeda078009cf3db74e90d5d42328/certificates/key_store/9b7d1e723af8233/download"
		Context(`Using mock server endpoint with timeout`, func() {
			BeforeEach(func() {
				testServer = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
					defer GinkgoRecover()

					// Verify the contents of the request
					Expect(req.URL.EscapedPath()).To(Equal(downloadKeyStoreCertificatePath))
					Expect(req.Method).To(Equal("GET"))

					Expect(req.Header["Accept-Language"]).ToNot(BeNil())
					Expect(req.Header["Accept-Language"][0]).To(Equal(fmt.Sprintf("%v", "en-US,en;q=0.5")))
					// Sleep a short time to support a timeout test
					time.Sleep(100 * time.Millisecond)

					// Set mock response
					res.Header().Set("Content-type", "application/octet-stream")
					res.WriteHeader(200)
					fmt.Fprintf(res, "%s", `This is a mock binary response.`)
				}))
			})
			It(`Invoke DownloadKeyStoreCertificate successfully with retries`, func() {
				mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1(&mqcloudv1.MqcloudV1Options{
					URL:            testServer.URL,
					Authenticator:  &core.NoAuthAuthenticator{},
					AcceptLanguage: core.StringPtr(acceptLanguage),
				})
				Expect(serviceErr).To(BeNil())
				Expect(mqcloudService).ToNot(BeNil())
				mqcloudService.EnableRetries(0, 0)

				// Construct an instance of the DownloadKeyStoreCertificateOptions model
				downloadKeyStoreCertificateOptionsModel := new(mqcloudv1.DownloadKeyStoreCertificateOptions)
				downloadKeyStoreCertificateOptionsModel.ServiceInstanceGuid = core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				downloadKeyStoreCertificateOptionsModel.QueueManagerID = core.StringPtr("b8e1aeda078009cf3db74e90d5d42328")
				downloadKeyStoreCertificateOptionsModel.CertificateID = core.StringPtr("9b7d1e723af8233")
				downloadKeyStoreCertificateOptionsModel.Headers = map[string]string{"x-custom-header": "x-custom-value"}

				// Invoke operation with a Context to test a timeout error
				ctx, cancelFunc := context.WithTimeout(context.Background(), 80*time.Millisecond)
				defer cancelFunc()
				_, _, operationErr := mqcloudService.DownloadKeyStoreCertificateWithContext(ctx, downloadKeyStoreCertificateOptionsModel)
				Expect(operationErr).ToNot(BeNil())
				Expect(operationErr.Error()).To(ContainSubstring("deadline exceeded"))

				// Disable retries and test again
				mqcloudService.DisableRetries()
				result, response, operationErr := mqcloudService.DownloadKeyStoreCertificate(downloadKeyStoreCertificateOptionsModel)
				Expect(operationErr).To(BeNil())
				Expect(response).ToNot(BeNil())
				Expect(result).ToNot(BeNil())

				// Re-test the timeout error with retries disabled
				ctx, cancelFunc2 := context.WithTimeout(context.Background(), 80*time.Millisecond)
				defer cancelFunc2()
				_, _, operationErr = mqcloudService.DownloadKeyStoreCertificateWithContext(ctx, downloadKeyStoreCertificateOptionsModel)
				Expect(operationErr).ToNot(BeNil())
				Expect(operationErr.Error()).To(ContainSubstring("deadline exceeded"))
			})
			AfterEach(func() {
				testServer.Close()
			})
		})
		Context(`Using mock server endpoint`, func() {
			BeforeEach(func() {
				testServer = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
					defer GinkgoRecover()

					// Verify the contents of the request
					Expect(req.URL.EscapedPath()).To(Equal(downloadKeyStoreCertificatePath))
					Expect(req.Method).To(Equal("GET"))

					Expect(req.Header["Accept-Language"]).ToNot(BeNil())
					Expect(req.Header["Accept-Language"][0]).To(Equal(fmt.Sprintf("%v", "en-US,en;q=0.5")))
					// Set mock response
					res.Header().Set("Content-type", "application/octet-stream")
					res.WriteHeader(200)
					fmt.Fprintf(res, "%s", `This is a mock binary response.`)
				}))
			})
			It(`Invoke DownloadKeyStoreCertificate successfully`, func() {
				mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1(&mqcloudv1.MqcloudV1Options{
					URL:            testServer.URL,
					Authenticator:  &core.NoAuthAuthenticator{},
					AcceptLanguage: core.StringPtr(acceptLanguage),
				})
				Expect(serviceErr).To(BeNil())
				Expect(mqcloudService).ToNot(BeNil())

				// Invoke operation with nil options model (negative test)
				result, response, operationErr := mqcloudService.DownloadKeyStoreCertificate(nil)
				Expect(operationErr).NotTo(BeNil())
				Expect(response).To(BeNil())
				Expect(result).To(BeNil())

				// Construct an instance of the DownloadKeyStoreCertificateOptions model
				downloadKeyStoreCertificateOptionsModel := new(mqcloudv1.DownloadKeyStoreCertificateOptions)
				downloadKeyStoreCertificateOptionsModel.ServiceInstanceGuid = core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				downloadKeyStoreCertificateOptionsModel.QueueManagerID = core.StringPtr("b8e1aeda078009cf3db74e90d5d42328")
				downloadKeyStoreCertificateOptionsModel.CertificateID = core.StringPtr("9b7d1e723af8233")
				downloadKeyStoreCertificateOptionsModel.Headers = map[string]string{"x-custom-header": "x-custom-value"}

				// Invoke operation with valid options model (positive test)
				result, response, operationErr = mqcloudService.DownloadKeyStoreCertificate(downloadKeyStoreCertificateOptionsModel)
				Expect(operationErr).To(BeNil())
				Expect(response).ToNot(BeNil())
				Expect(result).ToNot(BeNil())

			})
			It(`Invoke DownloadKeyStoreCertificate with error: Operation validation and request error`, func() {
				mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1(&mqcloudv1.MqcloudV1Options{
					URL:            testServer.URL,
					Authenticator:  &core.NoAuthAuthenticator{},
					AcceptLanguage: core.StringPtr(acceptLanguage),
				})
				Expect(serviceErr).To(BeNil())
				Expect(mqcloudService).ToNot(BeNil())

				// Construct an instance of the DownloadKeyStoreCertificateOptions model
				downloadKeyStoreCertificateOptionsModel := new(mqcloudv1.DownloadKeyStoreCertificateOptions)
				downloadKeyStoreCertificateOptionsModel.ServiceInstanceGuid = core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				downloadKeyStoreCertificateOptionsModel.QueueManagerID = core.StringPtr("b8e1aeda078009cf3db74e90d5d42328")
				downloadKeyStoreCertificateOptionsModel.CertificateID = core.StringPtr("9b7d1e723af8233")
				downloadKeyStoreCertificateOptionsModel.Headers = map[string]string{"x-custom-header": "x-custom-value"}
				// Invoke operation with empty URL (negative test)
				err := mqcloudService.SetServiceURL("")
				Expect(err).To(BeNil())
				result, response, operationErr := mqcloudService.DownloadKeyStoreCertificate(downloadKeyStoreCertificateOptionsModel)
				Expect(operationErr).ToNot(BeNil())
				Expect(operationErr.Error()).To(ContainSubstring(core.ERRORMSG_SERVICE_URL_MISSING))
				Expect(response).To(BeNil())
				Expect(result).To(BeNil())
				// Construct a second instance of the DownloadKeyStoreCertificateOptions model with no property values
				downloadKeyStoreCertificateOptionsModelNew := new(mqcloudv1.DownloadKeyStoreCertificateOptions)
				// Invoke operation with invalid model (negative test)
				result, response, operationErr = mqcloudService.DownloadKeyStoreCertificate(downloadKeyStoreCertificateOptionsModelNew)
				Expect(operationErr).ToNot(BeNil())
				Expect(response).To(BeNil())
				Expect(result).To(BeNil())
			})
			AfterEach(func() {
				testServer.Close()
			})
		})
		Context(`Using mock server endpoint with missing response body`, func() {
			BeforeEach(func() {
				testServer = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
					defer GinkgoRecover()

					// Set success status code with no respoonse body
					res.WriteHeader(200)
				}))
			})
			It(`Invoke DownloadKeyStoreCertificate successfully`, func() {
				mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1(&mqcloudv1.MqcloudV1Options{
					URL:            testServer.URL,
					Authenticator:  &core.NoAuthAuthenticator{},
					AcceptLanguage: core.StringPtr(acceptLanguage),
				})
				Expect(serviceErr).To(BeNil())
				Expect(mqcloudService).ToNot(BeNil())

				// Construct an instance of the DownloadKeyStoreCertificateOptions model
				downloadKeyStoreCertificateOptionsModel := new(mqcloudv1.DownloadKeyStoreCertificateOptions)
				downloadKeyStoreCertificateOptionsModel.ServiceInstanceGuid = core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				downloadKeyStoreCertificateOptionsModel.QueueManagerID = core.StringPtr("b8e1aeda078009cf3db74e90d5d42328")
				downloadKeyStoreCertificateOptionsModel.CertificateID = core.StringPtr("9b7d1e723af8233")
				downloadKeyStoreCertificateOptionsModel.Headers = map[string]string{"x-custom-header": "x-custom-value"}

				// Invoke operation
				result, response, operationErr := mqcloudService.DownloadKeyStoreCertificate(downloadKeyStoreCertificateOptionsModel)
				Expect(operationErr).To(BeNil())
				Expect(response).ToNot(BeNil())

				// Verify empty byte buffer.
				Expect(result).ToNot(BeNil())
				buffer, operationErr := io.ReadAll(result)
				Expect(operationErr).To(BeNil())
				Expect(buffer).ToNot(BeNil())
				Expect(len(buffer)).To(Equal(0))
			})
			AfterEach(func() {
				testServer.Close()
			})
		})
	})
	Describe(`GetCertificateAmsChannels(getCertificateAmsChannelsOptions *GetCertificateAmsChannelsOptions) - Operation response error`, func() {
		acceptLanguage := "en-US,en;q=0.5"
		getCertificateAmsChannelsPath := "/v1/a2b4d4bc-dadb-4637-bcec-9b7d1e723af8/queue_managers/b8e1aeda078009cf3db74e90d5d42328/certificates/key_store/9b7d1e723af8233/config/ams"
		Context(`Using mock server endpoint with invalid JSON response`, func() {
			BeforeEach(func() {
				testServer = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
					defer GinkgoRecover()

					// Verify the contents of the request
					Expect(req.URL.EscapedPath()).To(Equal(getCertificateAmsChannelsPath))
					Expect(req.Method).To(Equal("GET"))
					Expect(req.Header["Accept-Language"]).ToNot(BeNil())
					Expect(req.Header["Accept-Language"][0]).To(Equal(fmt.Sprintf("%v", "en-US,en;q=0.5")))
					res.Header().Set("Content-type", "application/json")
					res.WriteHeader(200)
					fmt.Fprint(res, `} this is not valid json {`)
				}))
			})
			It(`Invoke GetCertificateAmsChannels with error: Operation response processing error`, func() {
				mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1(&mqcloudv1.MqcloudV1Options{
					URL:            testServer.URL,
					Authenticator:  &core.NoAuthAuthenticator{},
					AcceptLanguage: core.StringPtr(acceptLanguage),
				})
				Expect(serviceErr).To(BeNil())
				Expect(mqcloudService).ToNot(BeNil())

				// Construct an instance of the GetCertificateAmsChannelsOptions model
				getCertificateAmsChannelsOptionsModel := new(mqcloudv1.GetCertificateAmsChannelsOptions)
				getCertificateAmsChannelsOptionsModel.QueueManagerID = core.StringPtr("b8e1aeda078009cf3db74e90d5d42328")
				getCertificateAmsChannelsOptionsModel.CertificateID = core.StringPtr("9b7d1e723af8233")
				getCertificateAmsChannelsOptionsModel.ServiceInstanceGuid = core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				getCertificateAmsChannelsOptionsModel.Headers = map[string]string{"x-custom-header": "x-custom-value"}
				// Expect response parsing to fail since we are receiving a text/plain response
				result, response, operationErr := mqcloudService.GetCertificateAmsChannels(getCertificateAmsChannelsOptionsModel)
				Expect(operationErr).ToNot(BeNil())
				Expect(response).ToNot(BeNil())
				Expect(result).To(BeNil())

				// Enable retries and test again
				mqcloudService.EnableRetries(0, 0)
				result, response, operationErr = mqcloudService.GetCertificateAmsChannels(getCertificateAmsChannelsOptionsModel)
				Expect(operationErr).ToNot(BeNil())
				Expect(response).ToNot(BeNil())
				Expect(result).To(BeNil())
			})
			AfterEach(func() {
				testServer.Close()
			})
		})
	})
	Describe(`GetCertificateAmsChannels(getCertificateAmsChannelsOptions *GetCertificateAmsChannelsOptions)`, func() {
		acceptLanguage := "en-US,en;q=0.5"
		getCertificateAmsChannelsPath := "/v1/a2b4d4bc-dadb-4637-bcec-9b7d1e723af8/queue_managers/b8e1aeda078009cf3db74e90d5d42328/certificates/key_store/9b7d1e723af8233/config/ams"
		Context(`Using mock server endpoint with timeout`, func() {
			BeforeEach(func() {
				testServer = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
					defer GinkgoRecover()

					// Verify the contents of the request
					Expect(req.URL.EscapedPath()).To(Equal(getCertificateAmsChannelsPath))
					Expect(req.Method).To(Equal("GET"))

					Expect(req.Header["Accept-Language"]).ToNot(BeNil())
					Expect(req.Header["Accept-Language"][0]).To(Equal(fmt.Sprintf("%v", "en-US,en;q=0.5")))
					// Sleep a short time to support a timeout test
					time.Sleep(100 * time.Millisecond)

					// Set mock response
					res.Header().Set("Content-type", "application/json")
					res.WriteHeader(200)
					fmt.Fprintf(res, "%s", `{"channels": [{"name": "Name"}]}`)
				}))
			})
			It(`Invoke GetCertificateAmsChannels successfully with retries`, func() {
				mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1(&mqcloudv1.MqcloudV1Options{
					URL:            testServer.URL,
					Authenticator:  &core.NoAuthAuthenticator{},
					AcceptLanguage: core.StringPtr(acceptLanguage),
				})
				Expect(serviceErr).To(BeNil())
				Expect(mqcloudService).ToNot(BeNil())
				mqcloudService.EnableRetries(0, 0)

				// Construct an instance of the GetCertificateAmsChannelsOptions model
				getCertificateAmsChannelsOptionsModel := new(mqcloudv1.GetCertificateAmsChannelsOptions)
				getCertificateAmsChannelsOptionsModel.QueueManagerID = core.StringPtr("b8e1aeda078009cf3db74e90d5d42328")
				getCertificateAmsChannelsOptionsModel.CertificateID = core.StringPtr("9b7d1e723af8233")
				getCertificateAmsChannelsOptionsModel.ServiceInstanceGuid = core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				getCertificateAmsChannelsOptionsModel.Headers = map[string]string{"x-custom-header": "x-custom-value"}

				// Invoke operation with a Context to test a timeout error
				ctx, cancelFunc := context.WithTimeout(context.Background(), 80*time.Millisecond)
				defer cancelFunc()
				_, _, operationErr := mqcloudService.GetCertificateAmsChannelsWithContext(ctx, getCertificateAmsChannelsOptionsModel)
				Expect(operationErr).ToNot(BeNil())
				Expect(operationErr.Error()).To(ContainSubstring("deadline exceeded"))

				// Disable retries and test again
				mqcloudService.DisableRetries()
				result, response, operationErr := mqcloudService.GetCertificateAmsChannels(getCertificateAmsChannelsOptionsModel)
				Expect(operationErr).To(BeNil())
				Expect(response).ToNot(BeNil())
				Expect(result).ToNot(BeNil())

				// Re-test the timeout error with retries disabled
				ctx, cancelFunc2 := context.WithTimeout(context.Background(), 80*time.Millisecond)
				defer cancelFunc2()
				_, _, operationErr = mqcloudService.GetCertificateAmsChannelsWithContext(ctx, getCertificateAmsChannelsOptionsModel)
				Expect(operationErr).ToNot(BeNil())
				Expect(operationErr.Error()).To(ContainSubstring("deadline exceeded"))
			})
			AfterEach(func() {
				testServer.Close()
			})
		})
		Context(`Using mock server endpoint`, func() {
			BeforeEach(func() {
				testServer = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
					defer GinkgoRecover()

					// Verify the contents of the request
					Expect(req.URL.EscapedPath()).To(Equal(getCertificateAmsChannelsPath))
					Expect(req.Method).To(Equal("GET"))

					Expect(req.Header["Accept-Language"]).ToNot(BeNil())
					Expect(req.Header["Accept-Language"][0]).To(Equal(fmt.Sprintf("%v", "en-US,en;q=0.5")))
					// Set mock response
					res.Header().Set("Content-type", "application/json")
					res.WriteHeader(200)
					fmt.Fprintf(res, "%s", `{"channels": [{"name": "Name"}]}`)
				}))
			})
			It(`Invoke GetCertificateAmsChannels successfully`, func() {
				mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1(&mqcloudv1.MqcloudV1Options{
					URL:            testServer.URL,
					Authenticator:  &core.NoAuthAuthenticator{},
					AcceptLanguage: core.StringPtr(acceptLanguage),
				})
				Expect(serviceErr).To(BeNil())
				Expect(mqcloudService).ToNot(BeNil())

				// Invoke operation with nil options model (negative test)
				result, response, operationErr := mqcloudService.GetCertificateAmsChannels(nil)
				Expect(operationErr).NotTo(BeNil())
				Expect(response).To(BeNil())
				Expect(result).To(BeNil())

				// Construct an instance of the GetCertificateAmsChannelsOptions model
				getCertificateAmsChannelsOptionsModel := new(mqcloudv1.GetCertificateAmsChannelsOptions)
				getCertificateAmsChannelsOptionsModel.QueueManagerID = core.StringPtr("b8e1aeda078009cf3db74e90d5d42328")
				getCertificateAmsChannelsOptionsModel.CertificateID = core.StringPtr("9b7d1e723af8233")
				getCertificateAmsChannelsOptionsModel.ServiceInstanceGuid = core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				getCertificateAmsChannelsOptionsModel.Headers = map[string]string{"x-custom-header": "x-custom-value"}

				// Invoke operation with valid options model (positive test)
				result, response, operationErr = mqcloudService.GetCertificateAmsChannels(getCertificateAmsChannelsOptionsModel)
				Expect(operationErr).To(BeNil())
				Expect(response).ToNot(BeNil())
				Expect(result).ToNot(BeNil())

			})
			It(`Invoke GetCertificateAmsChannels with error: Operation validation and request error`, func() {
				mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1(&mqcloudv1.MqcloudV1Options{
					URL:            testServer.URL,
					Authenticator:  &core.NoAuthAuthenticator{},
					AcceptLanguage: core.StringPtr(acceptLanguage),
				})
				Expect(serviceErr).To(BeNil())
				Expect(mqcloudService).ToNot(BeNil())

				// Construct an instance of the GetCertificateAmsChannelsOptions model
				getCertificateAmsChannelsOptionsModel := new(mqcloudv1.GetCertificateAmsChannelsOptions)
				getCertificateAmsChannelsOptionsModel.QueueManagerID = core.StringPtr("b8e1aeda078009cf3db74e90d5d42328")
				getCertificateAmsChannelsOptionsModel.CertificateID = core.StringPtr("9b7d1e723af8233")
				getCertificateAmsChannelsOptionsModel.ServiceInstanceGuid = core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				getCertificateAmsChannelsOptionsModel.Headers = map[string]string{"x-custom-header": "x-custom-value"}
				// Invoke operation with empty URL (negative test)
				err := mqcloudService.SetServiceURL("")
				Expect(err).To(BeNil())
				result, response, operationErr := mqcloudService.GetCertificateAmsChannels(getCertificateAmsChannelsOptionsModel)
				Expect(operationErr).ToNot(BeNil())
				Expect(operationErr.Error()).To(ContainSubstring(core.ERRORMSG_SERVICE_URL_MISSING))
				Expect(response).To(BeNil())
				Expect(result).To(BeNil())
				// Construct a second instance of the GetCertificateAmsChannelsOptions model with no property values
				getCertificateAmsChannelsOptionsModelNew := new(mqcloudv1.GetCertificateAmsChannelsOptions)
				// Invoke operation with invalid model (negative test)
				result, response, operationErr = mqcloudService.GetCertificateAmsChannels(getCertificateAmsChannelsOptionsModelNew)
				Expect(operationErr).ToNot(BeNil())
				Expect(response).To(BeNil())
				Expect(result).To(BeNil())
			})
			AfterEach(func() {
				testServer.Close()
			})
		})
		Context(`Using mock server endpoint with missing response body`, func() {
			BeforeEach(func() {
				testServer = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
					defer GinkgoRecover()

					// Set success status code with no respoonse body
					res.WriteHeader(200)
				}))
			})
			It(`Invoke GetCertificateAmsChannels successfully`, func() {
				mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1(&mqcloudv1.MqcloudV1Options{
					URL:            testServer.URL,
					Authenticator:  &core.NoAuthAuthenticator{},
					AcceptLanguage: core.StringPtr(acceptLanguage),
				})
				Expect(serviceErr).To(BeNil())
				Expect(mqcloudService).ToNot(BeNil())

				// Construct an instance of the GetCertificateAmsChannelsOptions model
				getCertificateAmsChannelsOptionsModel := new(mqcloudv1.GetCertificateAmsChannelsOptions)
				getCertificateAmsChannelsOptionsModel.QueueManagerID = core.StringPtr("b8e1aeda078009cf3db74e90d5d42328")
				getCertificateAmsChannelsOptionsModel.CertificateID = core.StringPtr("9b7d1e723af8233")
				getCertificateAmsChannelsOptionsModel.ServiceInstanceGuid = core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				getCertificateAmsChannelsOptionsModel.Headers = map[string]string{"x-custom-header": "x-custom-value"}

				// Invoke operation
				result, response, operationErr := mqcloudService.GetCertificateAmsChannels(getCertificateAmsChannelsOptionsModel)
				Expect(operationErr).To(BeNil())
				Expect(response).ToNot(BeNil())

				// Verify a nil result
				Expect(result).To(BeNil())
			})
			AfterEach(func() {
				testServer.Close()
			})
		})
	})
	Describe(`SetCertificateAmsChannels(setCertificateAmsChannelsOptions *SetCertificateAmsChannelsOptions) - Operation response error`, func() {
		acceptLanguage := "en-US,en;q=0.5"
		setCertificateAmsChannelsPath := "/v1/a2b4d4bc-dadb-4637-bcec-9b7d1e723af8/queue_managers/b8e1aeda078009cf3db74e90d5d42328/certificates/key_store/9b7d1e723af8233/config/ams"
		Context(`Using mock server endpoint with invalid JSON response`, func() {
			BeforeEach(func() {
				testServer = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
					defer GinkgoRecover()

					// Verify the contents of the request
					Expect(req.URL.EscapedPath()).To(Equal(setCertificateAmsChannelsPath))
					Expect(req.Method).To(Equal("PUT"))
					Expect(req.Header["Accept-Language"]).ToNot(BeNil())
					Expect(req.Header["Accept-Language"][0]).To(Equal(fmt.Sprintf("%v", "en-US,en;q=0.5")))
					res.Header().Set("Content-type", "application/json")
					res.WriteHeader(200)
					fmt.Fprint(res, `} this is not valid json {`)
				}))
			})
			It(`Invoke SetCertificateAmsChannels with error: Operation response processing error`, func() {
				mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1(&mqcloudv1.MqcloudV1Options{
					URL:            testServer.URL,
					Authenticator:  &core.NoAuthAuthenticator{},
					AcceptLanguage: core.StringPtr(acceptLanguage),
				})
				Expect(serviceErr).To(BeNil())
				Expect(mqcloudService).ToNot(BeNil())

				// Construct an instance of the ChannelDetails model
				channelDetailsModel := new(mqcloudv1.ChannelDetails)
				channelDetailsModel.Name = core.StringPtr("testString")

				// Construct an instance of the SetCertificateAmsChannelsOptions model
				setCertificateAmsChannelsOptionsModel := new(mqcloudv1.SetCertificateAmsChannelsOptions)
				setCertificateAmsChannelsOptionsModel.QueueManagerID = core.StringPtr("b8e1aeda078009cf3db74e90d5d42328")
				setCertificateAmsChannelsOptionsModel.CertificateID = core.StringPtr("9b7d1e723af8233")
				setCertificateAmsChannelsOptionsModel.ServiceInstanceGuid = core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				setCertificateAmsChannelsOptionsModel.Channels = []mqcloudv1.ChannelDetails{*channelDetailsModel}
				setCertificateAmsChannelsOptionsModel.UpdateStrategy = core.StringPtr("replace")
				setCertificateAmsChannelsOptionsModel.Headers = map[string]string{"x-custom-header": "x-custom-value"}
				// Expect response parsing to fail since we are receiving a text/plain response
				result, response, operationErr := mqcloudService.SetCertificateAmsChannels(setCertificateAmsChannelsOptionsModel)
				Expect(operationErr).ToNot(BeNil())
				Expect(response).ToNot(BeNil())
				Expect(result).To(BeNil())

				// Enable retries and test again
				mqcloudService.EnableRetries(0, 0)
				result, response, operationErr = mqcloudService.SetCertificateAmsChannels(setCertificateAmsChannelsOptionsModel)
				Expect(operationErr).ToNot(BeNil())
				Expect(response).ToNot(BeNil())
				Expect(result).To(BeNil())
			})
			AfterEach(func() {
				testServer.Close()
			})
		})
	})
	Describe(`SetCertificateAmsChannels(setCertificateAmsChannelsOptions *SetCertificateAmsChannelsOptions)`, func() {
		acceptLanguage := "en-US,en;q=0.5"
		setCertificateAmsChannelsPath := "/v1/a2b4d4bc-dadb-4637-bcec-9b7d1e723af8/queue_managers/b8e1aeda078009cf3db74e90d5d42328/certificates/key_store/9b7d1e723af8233/config/ams"
		Context(`Using mock server endpoint with timeout`, func() {
			BeforeEach(func() {
				testServer = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
					defer GinkgoRecover()

					// Verify the contents of the request
					Expect(req.URL.EscapedPath()).To(Equal(setCertificateAmsChannelsPath))
					Expect(req.Method).To(Equal("PUT"))

					// For gzip-disabled operation, verify Content-Encoding is not set.
					Expect(req.Header.Get("Content-Encoding")).To(BeEmpty())

					// If there is a body, then make sure we can read it
					bodyBuf := new(bytes.Buffer)
					if req.Header.Get("Content-Encoding") == "gzip" {
						body, err := core.NewGzipDecompressionReader(req.Body)
						Expect(err).To(BeNil())
						_, err = bodyBuf.ReadFrom(body)
						Expect(err).To(BeNil())
					} else {
						_, err := bodyBuf.ReadFrom(req.Body)
						Expect(err).To(BeNil())
					}
					fmt.Fprintf(GinkgoWriter, "  Request body: %s", bodyBuf.String())

					Expect(req.Header["Accept-Language"]).ToNot(BeNil())
					Expect(req.Header["Accept-Language"][0]).To(Equal(fmt.Sprintf("%v", "en-US,en;q=0.5")))
					// Sleep a short time to support a timeout test
					time.Sleep(100 * time.Millisecond)

					// Set mock response
					res.Header().Set("Content-type", "application/json")
					res.WriteHeader(200)
					fmt.Fprintf(res, "%s", `{"channels": [{"name": "Name"}]}`)
				}))
			})
			It(`Invoke SetCertificateAmsChannels successfully with retries`, func() {
				mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1(&mqcloudv1.MqcloudV1Options{
					URL:            testServer.URL,
					Authenticator:  &core.NoAuthAuthenticator{},
					AcceptLanguage: core.StringPtr(acceptLanguage),
				})
				Expect(serviceErr).To(BeNil())
				Expect(mqcloudService).ToNot(BeNil())
				mqcloudService.EnableRetries(0, 0)

				// Construct an instance of the ChannelDetails model
				channelDetailsModel := new(mqcloudv1.ChannelDetails)
				channelDetailsModel.Name = core.StringPtr("testString")

				// Construct an instance of the SetCertificateAmsChannelsOptions model
				setCertificateAmsChannelsOptionsModel := new(mqcloudv1.SetCertificateAmsChannelsOptions)
				setCertificateAmsChannelsOptionsModel.QueueManagerID = core.StringPtr("b8e1aeda078009cf3db74e90d5d42328")
				setCertificateAmsChannelsOptionsModel.CertificateID = core.StringPtr("9b7d1e723af8233")
				setCertificateAmsChannelsOptionsModel.ServiceInstanceGuid = core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				setCertificateAmsChannelsOptionsModel.Channels = []mqcloudv1.ChannelDetails{*channelDetailsModel}
				setCertificateAmsChannelsOptionsModel.UpdateStrategy = core.StringPtr("replace")
				setCertificateAmsChannelsOptionsModel.Headers = map[string]string{"x-custom-header": "x-custom-value"}

				// Invoke operation with a Context to test a timeout error
				ctx, cancelFunc := context.WithTimeout(context.Background(), 80*time.Millisecond)
				defer cancelFunc()
				_, _, operationErr := mqcloudService.SetCertificateAmsChannelsWithContext(ctx, setCertificateAmsChannelsOptionsModel)
				Expect(operationErr).ToNot(BeNil())
				Expect(operationErr.Error()).To(ContainSubstring("deadline exceeded"))

				// Disable retries and test again
				mqcloudService.DisableRetries()
				result, response, operationErr := mqcloudService.SetCertificateAmsChannels(setCertificateAmsChannelsOptionsModel)
				Expect(operationErr).To(BeNil())
				Expect(response).ToNot(BeNil())
				Expect(result).ToNot(BeNil())

				// Re-test the timeout error with retries disabled
				ctx, cancelFunc2 := context.WithTimeout(context.Background(), 80*time.Millisecond)
				defer cancelFunc2()
				_, _, operationErr = mqcloudService.SetCertificateAmsChannelsWithContext(ctx, setCertificateAmsChannelsOptionsModel)
				Expect(operationErr).ToNot(BeNil())
				Expect(operationErr.Error()).To(ContainSubstring("deadline exceeded"))
			})
			AfterEach(func() {
				testServer.Close()
			})
		})
		Context(`Using mock server endpoint`, func() {
			BeforeEach(func() {
				testServer = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
					defer GinkgoRecover()

					// Verify the contents of the request
					Expect(req.URL.EscapedPath()).To(Equal(setCertificateAmsChannelsPath))
					Expect(req.Method).To(Equal("PUT"))

					// For gzip-disabled operation, verify Content-Encoding is not set.
					Expect(req.Header.Get("Content-Encoding")).To(BeEmpty())

					// If there is a body, then make sure we can read it
					bodyBuf := new(bytes.Buffer)
					if req.Header.Get("Content-Encoding") == "gzip" {
						body, err := core.NewGzipDecompressionReader(req.Body)
						Expect(err).To(BeNil())
						_, err = bodyBuf.ReadFrom(body)
						Expect(err).To(BeNil())
					} else {
						_, err := bodyBuf.ReadFrom(req.Body)
						Expect(err).To(BeNil())
					}
					fmt.Fprintf(GinkgoWriter, "  Request body: %s", bodyBuf.String())

					Expect(req.Header["Accept-Language"]).ToNot(BeNil())
					Expect(req.Header["Accept-Language"][0]).To(Equal(fmt.Sprintf("%v", "en-US,en;q=0.5")))
					// Set mock response
					res.Header().Set("Content-type", "application/json")
					res.WriteHeader(200)
					fmt.Fprintf(res, "%s", `{"channels": [{"name": "Name"}]}`)
				}))
			})
			It(`Invoke SetCertificateAmsChannels successfully`, func() {
				mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1(&mqcloudv1.MqcloudV1Options{
					URL:            testServer.URL,
					Authenticator:  &core.NoAuthAuthenticator{},
					AcceptLanguage: core.StringPtr(acceptLanguage),
				})
				Expect(serviceErr).To(BeNil())
				Expect(mqcloudService).ToNot(BeNil())

				// Invoke operation with nil options model (negative test)
				result, response, operationErr := mqcloudService.SetCertificateAmsChannels(nil)
				Expect(operationErr).NotTo(BeNil())
				Expect(response).To(BeNil())
				Expect(result).To(BeNil())

				// Construct an instance of the ChannelDetails model
				channelDetailsModel := new(mqcloudv1.ChannelDetails)
				channelDetailsModel.Name = core.StringPtr("testString")

				// Construct an instance of the SetCertificateAmsChannelsOptions model
				setCertificateAmsChannelsOptionsModel := new(mqcloudv1.SetCertificateAmsChannelsOptions)
				setCertificateAmsChannelsOptionsModel.QueueManagerID = core.StringPtr("b8e1aeda078009cf3db74e90d5d42328")
				setCertificateAmsChannelsOptionsModel.CertificateID = core.StringPtr("9b7d1e723af8233")
				setCertificateAmsChannelsOptionsModel.ServiceInstanceGuid = core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				setCertificateAmsChannelsOptionsModel.Channels = []mqcloudv1.ChannelDetails{*channelDetailsModel}
				setCertificateAmsChannelsOptionsModel.UpdateStrategy = core.StringPtr("replace")
				setCertificateAmsChannelsOptionsModel.Headers = map[string]string{"x-custom-header": "x-custom-value"}

				// Invoke operation with valid options model (positive test)
				result, response, operationErr = mqcloudService.SetCertificateAmsChannels(setCertificateAmsChannelsOptionsModel)
				Expect(operationErr).To(BeNil())
				Expect(response).ToNot(BeNil())
				Expect(result).ToNot(BeNil())

			})
			It(`Invoke SetCertificateAmsChannels with error: Operation validation and request error`, func() {
				mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1(&mqcloudv1.MqcloudV1Options{
					URL:            testServer.URL,
					Authenticator:  &core.NoAuthAuthenticator{},
					AcceptLanguage: core.StringPtr(acceptLanguage),
				})
				Expect(serviceErr).To(BeNil())
				Expect(mqcloudService).ToNot(BeNil())

				// Construct an instance of the ChannelDetails model
				channelDetailsModel := new(mqcloudv1.ChannelDetails)
				channelDetailsModel.Name = core.StringPtr("testString")

				// Construct an instance of the SetCertificateAmsChannelsOptions model
				setCertificateAmsChannelsOptionsModel := new(mqcloudv1.SetCertificateAmsChannelsOptions)
				setCertificateAmsChannelsOptionsModel.QueueManagerID = core.StringPtr("b8e1aeda078009cf3db74e90d5d42328")
				setCertificateAmsChannelsOptionsModel.CertificateID = core.StringPtr("9b7d1e723af8233")
				setCertificateAmsChannelsOptionsModel.ServiceInstanceGuid = core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				setCertificateAmsChannelsOptionsModel.Channels = []mqcloudv1.ChannelDetails{*channelDetailsModel}
				setCertificateAmsChannelsOptionsModel.UpdateStrategy = core.StringPtr("replace")
				setCertificateAmsChannelsOptionsModel.Headers = map[string]string{"x-custom-header": "x-custom-value"}
				// Invoke operation with empty URL (negative test)
				err := mqcloudService.SetServiceURL("")
				Expect(err).To(BeNil())
				result, response, operationErr := mqcloudService.SetCertificateAmsChannels(setCertificateAmsChannelsOptionsModel)
				Expect(operationErr).ToNot(BeNil())
				Expect(operationErr.Error()).To(ContainSubstring(core.ERRORMSG_SERVICE_URL_MISSING))
				Expect(response).To(BeNil())
				Expect(result).To(BeNil())
				// Construct a second instance of the SetCertificateAmsChannelsOptions model with no property values
				setCertificateAmsChannelsOptionsModelNew := new(mqcloudv1.SetCertificateAmsChannelsOptions)
				// Invoke operation with invalid model (negative test)
				result, response, operationErr = mqcloudService.SetCertificateAmsChannels(setCertificateAmsChannelsOptionsModelNew)
				Expect(operationErr).ToNot(BeNil())
				Expect(response).To(BeNil())
				Expect(result).To(BeNil())
			})
			AfterEach(func() {
				testServer.Close()
			})
		})
		Context(`Using mock server endpoint with missing response body`, func() {
			BeforeEach(func() {
				testServer = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
					defer GinkgoRecover()

					// Set success status code with no respoonse body
					res.WriteHeader(200)
				}))
			})
			It(`Invoke SetCertificateAmsChannels successfully`, func() {
				mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1(&mqcloudv1.MqcloudV1Options{
					URL:            testServer.URL,
					Authenticator:  &core.NoAuthAuthenticator{},
					AcceptLanguage: core.StringPtr(acceptLanguage),
				})
				Expect(serviceErr).To(BeNil())
				Expect(mqcloudService).ToNot(BeNil())

				// Construct an instance of the ChannelDetails model
				channelDetailsModel := new(mqcloudv1.ChannelDetails)
				channelDetailsModel.Name = core.StringPtr("testString")

				// Construct an instance of the SetCertificateAmsChannelsOptions model
				setCertificateAmsChannelsOptionsModel := new(mqcloudv1.SetCertificateAmsChannelsOptions)
				setCertificateAmsChannelsOptionsModel.QueueManagerID = core.StringPtr("b8e1aeda078009cf3db74e90d5d42328")
				setCertificateAmsChannelsOptionsModel.CertificateID = core.StringPtr("9b7d1e723af8233")
				setCertificateAmsChannelsOptionsModel.ServiceInstanceGuid = core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				setCertificateAmsChannelsOptionsModel.Channels = []mqcloudv1.ChannelDetails{*channelDetailsModel}
				setCertificateAmsChannelsOptionsModel.UpdateStrategy = core.StringPtr("replace")
				setCertificateAmsChannelsOptionsModel.Headers = map[string]string{"x-custom-header": "x-custom-value"}

				// Invoke operation
				result, response, operationErr := mqcloudService.SetCertificateAmsChannels(setCertificateAmsChannelsOptionsModel)
				Expect(operationErr).To(BeNil())
				Expect(response).ToNot(BeNil())

				// Verify a nil result
				Expect(result).To(BeNil())
			})
			AfterEach(func() {
				testServer.Close()
			})
		})
	})
	Describe(`CreateVirtualPrivateEndpointGateway(createVirtualPrivateEndpointGatewayOptions *CreateVirtualPrivateEndpointGatewayOptions) - Operation response error`, func() {
		acceptLanguage := "en-US,en;q=0.5"
		createVirtualPrivateEndpointGatewayPath := "/v1/a2b4d4bc-dadb-4637-bcec-9b7d1e723af8/virtual_private_endpoint_gateway"
		Context(`Using mock server endpoint with invalid JSON response`, func() {
			BeforeEach(func() {
				testServer = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
					defer GinkgoRecover()

					// Verify the contents of the request
					Expect(req.URL.EscapedPath()).To(Equal(createVirtualPrivateEndpointGatewayPath))
					Expect(req.Method).To(Equal("POST"))
					Expect(req.Header["Accept-Language"]).ToNot(BeNil())
					Expect(req.Header["Accept-Language"][0]).To(Equal(fmt.Sprintf("%v", "en-US,en;q=0.5")))
					Expect(req.Header["Trusted-Profile"]).ToNot(BeNil())
					Expect(req.Header["Trusted-Profile"][0]).To(Equal(fmt.Sprintf("%v", "crn:v1:bluemix:public:iam-identity::a/5d5ff2a9001c4055ab1408e9bf185f48::profile:Profile-1c0a8609-ca25-4ad2-a09b-aea472d34afc")))
					res.Header().Set("Content-type", "application/json")
					res.WriteHeader(201)
					fmt.Fprint(res, `} this is not valid json {`)
				}))
			})
			It(`Invoke CreateVirtualPrivateEndpointGateway with error: Operation response processing error`, func() {
				mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1(&mqcloudv1.MqcloudV1Options{
					URL:            testServer.URL,
					Authenticator:  &core.NoAuthAuthenticator{},
					AcceptLanguage: core.StringPtr(acceptLanguage),
				})
				Expect(serviceErr).To(BeNil())
				Expect(mqcloudService).ToNot(BeNil())

				// Construct an instance of the CreateVirtualPrivateEndpointGatewayOptions model
				createVirtualPrivateEndpointGatewayOptionsModel := new(mqcloudv1.CreateVirtualPrivateEndpointGatewayOptions)
				createVirtualPrivateEndpointGatewayOptionsModel.ServiceInstanceGuid = core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				createVirtualPrivateEndpointGatewayOptionsModel.Name = core.StringPtr("vpe_gateway1-to-vpe_gateway2")
				createVirtualPrivateEndpointGatewayOptionsModel.TargetCrn = core.StringPtr("crn:v1:bluemix:public:mqcloud:eu-de:::endpoint:qm1.private.eu-de.mq2.test.appdomain.cloud")
				createVirtualPrivateEndpointGatewayOptionsModel.TrustedProfile = core.StringPtr("crn:v1:bluemix:public:iam-identity::a/5d5ff2a9001c4055ab1408e9bf185f48::profile:Profile-1c0a8609-ca25-4ad2-a09b-aea472d34afc")
				createVirtualPrivateEndpointGatewayOptionsModel.Headers = map[string]string{"x-custom-header": "x-custom-value"}
				// Expect response parsing to fail since we are receiving a text/plain response
				result, response, operationErr := mqcloudService.CreateVirtualPrivateEndpointGateway(createVirtualPrivateEndpointGatewayOptionsModel)
				Expect(operationErr).ToNot(BeNil())
				Expect(response).ToNot(BeNil())
				Expect(result).To(BeNil())

				// Enable retries and test again
				mqcloudService.EnableRetries(0, 0)
				result, response, operationErr = mqcloudService.CreateVirtualPrivateEndpointGateway(createVirtualPrivateEndpointGatewayOptionsModel)
				Expect(operationErr).ToNot(BeNil())
				Expect(response).ToNot(BeNil())
				Expect(result).To(BeNil())
			})
			AfterEach(func() {
				testServer.Close()
			})
		})
	})
	Describe(`CreateVirtualPrivateEndpointGateway(createVirtualPrivateEndpointGatewayOptions *CreateVirtualPrivateEndpointGatewayOptions)`, func() {
		acceptLanguage := "en-US,en;q=0.5"
		createVirtualPrivateEndpointGatewayPath := "/v1/a2b4d4bc-dadb-4637-bcec-9b7d1e723af8/virtual_private_endpoint_gateway"
		Context(`Using mock server endpoint with timeout`, func() {
			BeforeEach(func() {
				testServer = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
					defer GinkgoRecover()

					// Verify the contents of the request
					Expect(req.URL.EscapedPath()).To(Equal(createVirtualPrivateEndpointGatewayPath))
					Expect(req.Method).To(Equal("POST"))

					// For gzip-disabled operation, verify Content-Encoding is not set.
					Expect(req.Header.Get("Content-Encoding")).To(BeEmpty())

					// If there is a body, then make sure we can read it
					bodyBuf := new(bytes.Buffer)
					if req.Header.Get("Content-Encoding") == "gzip" {
						body, err := core.NewGzipDecompressionReader(req.Body)
						Expect(err).To(BeNil())
						_, err = bodyBuf.ReadFrom(body)
						Expect(err).To(BeNil())
					} else {
						_, err := bodyBuf.ReadFrom(req.Body)
						Expect(err).To(BeNil())
					}
					fmt.Fprintf(GinkgoWriter, "  Request body: %s", bodyBuf.String())

					Expect(req.Header["Accept-Language"]).ToNot(BeNil())
					Expect(req.Header["Accept-Language"][0]).To(Equal(fmt.Sprintf("%v", "en-US,en;q=0.5")))
					Expect(req.Header["Trusted-Profile"]).ToNot(BeNil())
					Expect(req.Header["Trusted-Profile"][0]).To(Equal(fmt.Sprintf("%v", "crn:v1:bluemix:public:iam-identity::a/5d5ff2a9001c4055ab1408e9bf185f48::profile:Profile-1c0a8609-ca25-4ad2-a09b-aea472d34afc")))
					// Sleep a short time to support a timeout test
					time.Sleep(100 * time.Millisecond)

					// Set mock response
					res.Header().Set("Content-type", "application/json")
					res.WriteHeader(201)
					fmt.Fprintf(res, "%s", `{"href": "Href", "id": "ID", "name": "Name", "target_crn": "TargetCrn", "status": "Status"}`)
				}))
			})
			It(`Invoke CreateVirtualPrivateEndpointGateway successfully with retries`, func() {
				mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1(&mqcloudv1.MqcloudV1Options{
					URL:            testServer.URL,
					Authenticator:  &core.NoAuthAuthenticator{},
					AcceptLanguage: core.StringPtr(acceptLanguage),
				})
				Expect(serviceErr).To(BeNil())
				Expect(mqcloudService).ToNot(BeNil())
				mqcloudService.EnableRetries(0, 0)

				// Construct an instance of the CreateVirtualPrivateEndpointGatewayOptions model
				createVirtualPrivateEndpointGatewayOptionsModel := new(mqcloudv1.CreateVirtualPrivateEndpointGatewayOptions)
				createVirtualPrivateEndpointGatewayOptionsModel.ServiceInstanceGuid = core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				createVirtualPrivateEndpointGatewayOptionsModel.Name = core.StringPtr("vpe_gateway1-to-vpe_gateway2")
				createVirtualPrivateEndpointGatewayOptionsModel.TargetCrn = core.StringPtr("crn:v1:bluemix:public:mqcloud:eu-de:::endpoint:qm1.private.eu-de.mq2.test.appdomain.cloud")
				createVirtualPrivateEndpointGatewayOptionsModel.TrustedProfile = core.StringPtr("crn:v1:bluemix:public:iam-identity::a/5d5ff2a9001c4055ab1408e9bf185f48::profile:Profile-1c0a8609-ca25-4ad2-a09b-aea472d34afc")
				createVirtualPrivateEndpointGatewayOptionsModel.Headers = map[string]string{"x-custom-header": "x-custom-value"}

				// Invoke operation with a Context to test a timeout error
				ctx, cancelFunc := context.WithTimeout(context.Background(), 80*time.Millisecond)
				defer cancelFunc()
				_, _, operationErr := mqcloudService.CreateVirtualPrivateEndpointGatewayWithContext(ctx, createVirtualPrivateEndpointGatewayOptionsModel)
				Expect(operationErr).ToNot(BeNil())
				Expect(operationErr.Error()).To(ContainSubstring("deadline exceeded"))

				// Disable retries and test again
				mqcloudService.DisableRetries()
				result, response, operationErr := mqcloudService.CreateVirtualPrivateEndpointGateway(createVirtualPrivateEndpointGatewayOptionsModel)
				Expect(operationErr).To(BeNil())
				Expect(response).ToNot(BeNil())
				Expect(result).ToNot(BeNil())

				// Re-test the timeout error with retries disabled
				ctx, cancelFunc2 := context.WithTimeout(context.Background(), 80*time.Millisecond)
				defer cancelFunc2()
				_, _, operationErr = mqcloudService.CreateVirtualPrivateEndpointGatewayWithContext(ctx, createVirtualPrivateEndpointGatewayOptionsModel)
				Expect(operationErr).ToNot(BeNil())
				Expect(operationErr.Error()).To(ContainSubstring("deadline exceeded"))
			})
			AfterEach(func() {
				testServer.Close()
			})
		})
		Context(`Using mock server endpoint`, func() {
			BeforeEach(func() {
				testServer = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
					defer GinkgoRecover()

					// Verify the contents of the request
					Expect(req.URL.EscapedPath()).To(Equal(createVirtualPrivateEndpointGatewayPath))
					Expect(req.Method).To(Equal("POST"))

					// For gzip-disabled operation, verify Content-Encoding is not set.
					Expect(req.Header.Get("Content-Encoding")).To(BeEmpty())

					// If there is a body, then make sure we can read it
					bodyBuf := new(bytes.Buffer)
					if req.Header.Get("Content-Encoding") == "gzip" {
						body, err := core.NewGzipDecompressionReader(req.Body)
						Expect(err).To(BeNil())
						_, err = bodyBuf.ReadFrom(body)
						Expect(err).To(BeNil())
					} else {
						_, err := bodyBuf.ReadFrom(req.Body)
						Expect(err).To(BeNil())
					}
					fmt.Fprintf(GinkgoWriter, "  Request body: %s", bodyBuf.String())

					Expect(req.Header["Accept-Language"]).ToNot(BeNil())
					Expect(req.Header["Accept-Language"][0]).To(Equal(fmt.Sprintf("%v", "en-US,en;q=0.5")))
					Expect(req.Header["Trusted-Profile"]).ToNot(BeNil())
					Expect(req.Header["Trusted-Profile"][0]).To(Equal(fmt.Sprintf("%v", "crn:v1:bluemix:public:iam-identity::a/5d5ff2a9001c4055ab1408e9bf185f48::profile:Profile-1c0a8609-ca25-4ad2-a09b-aea472d34afc")))
					// Set mock response
					res.Header().Set("Content-type", "application/json")
					res.WriteHeader(201)
					fmt.Fprintf(res, "%s", `{"href": "Href", "id": "ID", "name": "Name", "target_crn": "TargetCrn", "status": "Status"}`)
				}))
			})
			It(`Invoke CreateVirtualPrivateEndpointGateway successfully`, func() {
				mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1(&mqcloudv1.MqcloudV1Options{
					URL:            testServer.URL,
					Authenticator:  &core.NoAuthAuthenticator{},
					AcceptLanguage: core.StringPtr(acceptLanguage),
				})
				Expect(serviceErr).To(BeNil())
				Expect(mqcloudService).ToNot(BeNil())

				// Invoke operation with nil options model (negative test)
				result, response, operationErr := mqcloudService.CreateVirtualPrivateEndpointGateway(nil)
				Expect(operationErr).NotTo(BeNil())
				Expect(response).To(BeNil())
				Expect(result).To(BeNil())

				// Construct an instance of the CreateVirtualPrivateEndpointGatewayOptions model
				createVirtualPrivateEndpointGatewayOptionsModel := new(mqcloudv1.CreateVirtualPrivateEndpointGatewayOptions)
				createVirtualPrivateEndpointGatewayOptionsModel.ServiceInstanceGuid = core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				createVirtualPrivateEndpointGatewayOptionsModel.Name = core.StringPtr("vpe_gateway1-to-vpe_gateway2")
				createVirtualPrivateEndpointGatewayOptionsModel.TargetCrn = core.StringPtr("crn:v1:bluemix:public:mqcloud:eu-de:::endpoint:qm1.private.eu-de.mq2.test.appdomain.cloud")
				createVirtualPrivateEndpointGatewayOptionsModel.TrustedProfile = core.StringPtr("crn:v1:bluemix:public:iam-identity::a/5d5ff2a9001c4055ab1408e9bf185f48::profile:Profile-1c0a8609-ca25-4ad2-a09b-aea472d34afc")
				createVirtualPrivateEndpointGatewayOptionsModel.Headers = map[string]string{"x-custom-header": "x-custom-value"}

				// Invoke operation with valid options model (positive test)
				result, response, operationErr = mqcloudService.CreateVirtualPrivateEndpointGateway(createVirtualPrivateEndpointGatewayOptionsModel)
				Expect(operationErr).To(BeNil())
				Expect(response).ToNot(BeNil())
				Expect(result).ToNot(BeNil())

			})
			It(`Invoke CreateVirtualPrivateEndpointGateway with error: Operation validation and request error`, func() {
				mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1(&mqcloudv1.MqcloudV1Options{
					URL:            testServer.URL,
					Authenticator:  &core.NoAuthAuthenticator{},
					AcceptLanguage: core.StringPtr(acceptLanguage),
				})
				Expect(serviceErr).To(BeNil())
				Expect(mqcloudService).ToNot(BeNil())

				// Construct an instance of the CreateVirtualPrivateEndpointGatewayOptions model
				createVirtualPrivateEndpointGatewayOptionsModel := new(mqcloudv1.CreateVirtualPrivateEndpointGatewayOptions)
				createVirtualPrivateEndpointGatewayOptionsModel.ServiceInstanceGuid = core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				createVirtualPrivateEndpointGatewayOptionsModel.Name = core.StringPtr("vpe_gateway1-to-vpe_gateway2")
				createVirtualPrivateEndpointGatewayOptionsModel.TargetCrn = core.StringPtr("crn:v1:bluemix:public:mqcloud:eu-de:::endpoint:qm1.private.eu-de.mq2.test.appdomain.cloud")
				createVirtualPrivateEndpointGatewayOptionsModel.TrustedProfile = core.StringPtr("crn:v1:bluemix:public:iam-identity::a/5d5ff2a9001c4055ab1408e9bf185f48::profile:Profile-1c0a8609-ca25-4ad2-a09b-aea472d34afc")
				createVirtualPrivateEndpointGatewayOptionsModel.Headers = map[string]string{"x-custom-header": "x-custom-value"}
				// Invoke operation with empty URL (negative test)
				err := mqcloudService.SetServiceURL("")
				Expect(err).To(BeNil())
				result, response, operationErr := mqcloudService.CreateVirtualPrivateEndpointGateway(createVirtualPrivateEndpointGatewayOptionsModel)
				Expect(operationErr).ToNot(BeNil())
				Expect(operationErr.Error()).To(ContainSubstring(core.ERRORMSG_SERVICE_URL_MISSING))
				Expect(response).To(BeNil())
				Expect(result).To(BeNil())
				// Construct a second instance of the CreateVirtualPrivateEndpointGatewayOptions model with no property values
				createVirtualPrivateEndpointGatewayOptionsModelNew := new(mqcloudv1.CreateVirtualPrivateEndpointGatewayOptions)
				// Invoke operation with invalid model (negative test)
				result, response, operationErr = mqcloudService.CreateVirtualPrivateEndpointGateway(createVirtualPrivateEndpointGatewayOptionsModelNew)
				Expect(operationErr).ToNot(BeNil())
				Expect(response).To(BeNil())
				Expect(result).To(BeNil())
			})
			AfterEach(func() {
				testServer.Close()
			})
		})
		Context(`Using mock server endpoint with missing response body`, func() {
			BeforeEach(func() {
				testServer = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
					defer GinkgoRecover()

					// Set success status code with no respoonse body
					res.WriteHeader(201)
				}))
			})
			It(`Invoke CreateVirtualPrivateEndpointGateway successfully`, func() {
				mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1(&mqcloudv1.MqcloudV1Options{
					URL:            testServer.URL,
					Authenticator:  &core.NoAuthAuthenticator{},
					AcceptLanguage: core.StringPtr(acceptLanguage),
				})
				Expect(serviceErr).To(BeNil())
				Expect(mqcloudService).ToNot(BeNil())

				// Construct an instance of the CreateVirtualPrivateEndpointGatewayOptions model
				createVirtualPrivateEndpointGatewayOptionsModel := new(mqcloudv1.CreateVirtualPrivateEndpointGatewayOptions)
				createVirtualPrivateEndpointGatewayOptionsModel.ServiceInstanceGuid = core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				createVirtualPrivateEndpointGatewayOptionsModel.Name = core.StringPtr("vpe_gateway1-to-vpe_gateway2")
				createVirtualPrivateEndpointGatewayOptionsModel.TargetCrn = core.StringPtr("crn:v1:bluemix:public:mqcloud:eu-de:::endpoint:qm1.private.eu-de.mq2.test.appdomain.cloud")
				createVirtualPrivateEndpointGatewayOptionsModel.TrustedProfile = core.StringPtr("crn:v1:bluemix:public:iam-identity::a/5d5ff2a9001c4055ab1408e9bf185f48::profile:Profile-1c0a8609-ca25-4ad2-a09b-aea472d34afc")
				createVirtualPrivateEndpointGatewayOptionsModel.Headers = map[string]string{"x-custom-header": "x-custom-value"}

				// Invoke operation
				result, response, operationErr := mqcloudService.CreateVirtualPrivateEndpointGateway(createVirtualPrivateEndpointGatewayOptionsModel)
				Expect(operationErr).To(BeNil())
				Expect(response).ToNot(BeNil())

				// Verify a nil result
				Expect(result).To(BeNil())
			})
			AfterEach(func() {
				testServer.Close()
			})
		})
	})
	Describe(`ListVirtualPrivateEndpointGateways(listVirtualPrivateEndpointGatewaysOptions *ListVirtualPrivateEndpointGatewaysOptions) - Operation response error`, func() {
		acceptLanguage := "en-US,en;q=0.5"
		listVirtualPrivateEndpointGatewaysPath := "/v1/a2b4d4bc-dadb-4637-bcec-9b7d1e723af8/virtual_private_endpoint_gateway"
		Context(`Using mock server endpoint with invalid JSON response`, func() {
			BeforeEach(func() {
				testServer = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
					defer GinkgoRecover()

					// Verify the contents of the request
					Expect(req.URL.EscapedPath()).To(Equal(listVirtualPrivateEndpointGatewaysPath))
					Expect(req.Method).To(Equal("GET"))
					Expect(req.Header["Accept-Language"]).ToNot(BeNil())
					Expect(req.Header["Accept-Language"][0]).To(Equal(fmt.Sprintf("%v", "en-US,en;q=0.5")))
					Expect(req.Header["Trusted-Profile"]).ToNot(BeNil())
					Expect(req.Header["Trusted-Profile"][0]).To(Equal(fmt.Sprintf("%v", "crn:v1:bluemix:public:iam-identity::a/5d5ff2a9001c4055ab1408e9bf185f48::profile:Profile-1c0a8609-ca25-4ad2-a09b-aea472d34afc")))
					Expect(req.URL.Query()["start"]).To(Equal([]string{"r010-ebab3c08-c9a8-40c4-8869-61c09ddf7b44"}))
					Expect(req.URL.Query()["limit"]).To(Equal([]string{fmt.Sprint(int64(10))}))
					res.Header().Set("Content-type", "application/json")
					res.WriteHeader(200)
					fmt.Fprint(res, `} this is not valid json {`)
				}))
			})
			It(`Invoke ListVirtualPrivateEndpointGateways with error: Operation response processing error`, func() {
				mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1(&mqcloudv1.MqcloudV1Options{
					URL:            testServer.URL,
					Authenticator:  &core.NoAuthAuthenticator{},
					AcceptLanguage: core.StringPtr(acceptLanguage),
				})
				Expect(serviceErr).To(BeNil())
				Expect(mqcloudService).ToNot(BeNil())

				// Construct an instance of the ListVirtualPrivateEndpointGatewaysOptions model
				listVirtualPrivateEndpointGatewaysOptionsModel := new(mqcloudv1.ListVirtualPrivateEndpointGatewaysOptions)
				listVirtualPrivateEndpointGatewaysOptionsModel.ServiceInstanceGuid = core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				listVirtualPrivateEndpointGatewaysOptionsModel.TrustedProfile = core.StringPtr("crn:v1:bluemix:public:iam-identity::a/5d5ff2a9001c4055ab1408e9bf185f48::profile:Profile-1c0a8609-ca25-4ad2-a09b-aea472d34afc")
				listVirtualPrivateEndpointGatewaysOptionsModel.Start = core.StringPtr("r010-ebab3c08-c9a8-40c4-8869-61c09ddf7b44")
				listVirtualPrivateEndpointGatewaysOptionsModel.Limit = core.Int64Ptr(int64(10))
				listVirtualPrivateEndpointGatewaysOptionsModel.Headers = map[string]string{"x-custom-header": "x-custom-value"}
				// Expect response parsing to fail since we are receiving a text/plain response
				result, response, operationErr := mqcloudService.ListVirtualPrivateEndpointGateways(listVirtualPrivateEndpointGatewaysOptionsModel)
				Expect(operationErr).ToNot(BeNil())
				Expect(response).ToNot(BeNil())
				Expect(result).To(BeNil())

				// Enable retries and test again
				mqcloudService.EnableRetries(0, 0)
				result, response, operationErr = mqcloudService.ListVirtualPrivateEndpointGateways(listVirtualPrivateEndpointGatewaysOptionsModel)
				Expect(operationErr).ToNot(BeNil())
				Expect(response).ToNot(BeNil())
				Expect(result).To(BeNil())
			})
			AfterEach(func() {
				testServer.Close()
			})
		})
	})
	Describe(`ListVirtualPrivateEndpointGateways(listVirtualPrivateEndpointGatewaysOptions *ListVirtualPrivateEndpointGatewaysOptions)`, func() {
		acceptLanguage := "en-US,en;q=0.5"
		listVirtualPrivateEndpointGatewaysPath := "/v1/a2b4d4bc-dadb-4637-bcec-9b7d1e723af8/virtual_private_endpoint_gateway"
		Context(`Using mock server endpoint with timeout`, func() {
			BeforeEach(func() {
				testServer = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
					defer GinkgoRecover()

					// Verify the contents of the request
					Expect(req.URL.EscapedPath()).To(Equal(listVirtualPrivateEndpointGatewaysPath))
					Expect(req.Method).To(Equal("GET"))

					Expect(req.Header["Accept-Language"]).ToNot(BeNil())
					Expect(req.Header["Accept-Language"][0]).To(Equal(fmt.Sprintf("%v", "en-US,en;q=0.5")))
					Expect(req.Header["Trusted-Profile"]).ToNot(BeNil())
					Expect(req.Header["Trusted-Profile"][0]).To(Equal(fmt.Sprintf("%v", "crn:v1:bluemix:public:iam-identity::a/5d5ff2a9001c4055ab1408e9bf185f48::profile:Profile-1c0a8609-ca25-4ad2-a09b-aea472d34afc")))
					Expect(req.URL.Query()["start"]).To(Equal([]string{"r010-ebab3c08-c9a8-40c4-8869-61c09ddf7b44"}))
					Expect(req.URL.Query()["limit"]).To(Equal([]string{fmt.Sprint(int64(10))}))
					// Sleep a short time to support a timeout test
					time.Sleep(100 * time.Millisecond)

					// Set mock response
					res.Header().Set("Content-type", "application/json")
					res.WriteHeader(200)
					fmt.Fprintf(res, "%s", `{"limit": 50, "first": {"href": "Href"}, "next": {"href": "Href"}, "virtual_private_endpoint_gateways": [{"href": "Href", "id": "ID", "name": "Name", "target_crn": "TargetCrn", "status": "Status"}]}`)
				}))
			})
			It(`Invoke ListVirtualPrivateEndpointGateways successfully with retries`, func() {
				mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1(&mqcloudv1.MqcloudV1Options{
					URL:            testServer.URL,
					Authenticator:  &core.NoAuthAuthenticator{},
					AcceptLanguage: core.StringPtr(acceptLanguage),
				})
				Expect(serviceErr).To(BeNil())
				Expect(mqcloudService).ToNot(BeNil())
				mqcloudService.EnableRetries(0, 0)

				// Construct an instance of the ListVirtualPrivateEndpointGatewaysOptions model
				listVirtualPrivateEndpointGatewaysOptionsModel := new(mqcloudv1.ListVirtualPrivateEndpointGatewaysOptions)
				listVirtualPrivateEndpointGatewaysOptionsModel.ServiceInstanceGuid = core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				listVirtualPrivateEndpointGatewaysOptionsModel.TrustedProfile = core.StringPtr("crn:v1:bluemix:public:iam-identity::a/5d5ff2a9001c4055ab1408e9bf185f48::profile:Profile-1c0a8609-ca25-4ad2-a09b-aea472d34afc")
				listVirtualPrivateEndpointGatewaysOptionsModel.Start = core.StringPtr("r010-ebab3c08-c9a8-40c4-8869-61c09ddf7b44")
				listVirtualPrivateEndpointGatewaysOptionsModel.Limit = core.Int64Ptr(int64(10))
				listVirtualPrivateEndpointGatewaysOptionsModel.Headers = map[string]string{"x-custom-header": "x-custom-value"}

				// Invoke operation with a Context to test a timeout error
				ctx, cancelFunc := context.WithTimeout(context.Background(), 80*time.Millisecond)
				defer cancelFunc()
				_, _, operationErr := mqcloudService.ListVirtualPrivateEndpointGatewaysWithContext(ctx, listVirtualPrivateEndpointGatewaysOptionsModel)
				Expect(operationErr).ToNot(BeNil())
				Expect(operationErr.Error()).To(ContainSubstring("deadline exceeded"))

				// Disable retries and test again
				mqcloudService.DisableRetries()
				result, response, operationErr := mqcloudService.ListVirtualPrivateEndpointGateways(listVirtualPrivateEndpointGatewaysOptionsModel)
				Expect(operationErr).To(BeNil())
				Expect(response).ToNot(BeNil())
				Expect(result).ToNot(BeNil())

				// Re-test the timeout error with retries disabled
				ctx, cancelFunc2 := context.WithTimeout(context.Background(), 80*time.Millisecond)
				defer cancelFunc2()
				_, _, operationErr = mqcloudService.ListVirtualPrivateEndpointGatewaysWithContext(ctx, listVirtualPrivateEndpointGatewaysOptionsModel)
				Expect(operationErr).ToNot(BeNil())
				Expect(operationErr.Error()).To(ContainSubstring("deadline exceeded"))
			})
			AfterEach(func() {
				testServer.Close()
			})
		})
		Context(`Using mock server endpoint`, func() {
			BeforeEach(func() {
				testServer = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
					defer GinkgoRecover()

					// Verify the contents of the request
					Expect(req.URL.EscapedPath()).To(Equal(listVirtualPrivateEndpointGatewaysPath))
					Expect(req.Method).To(Equal("GET"))

					Expect(req.Header["Accept-Language"]).ToNot(BeNil())
					Expect(req.Header["Accept-Language"][0]).To(Equal(fmt.Sprintf("%v", "en-US,en;q=0.5")))
					Expect(req.Header["Trusted-Profile"]).ToNot(BeNil())
					Expect(req.Header["Trusted-Profile"][0]).To(Equal(fmt.Sprintf("%v", "crn:v1:bluemix:public:iam-identity::a/5d5ff2a9001c4055ab1408e9bf185f48::profile:Profile-1c0a8609-ca25-4ad2-a09b-aea472d34afc")))
					Expect(req.URL.Query()["start"]).To(Equal([]string{"r010-ebab3c08-c9a8-40c4-8869-61c09ddf7b44"}))
					Expect(req.URL.Query()["limit"]).To(Equal([]string{fmt.Sprint(int64(10))}))
					// Set mock response
					res.Header().Set("Content-type", "application/json")
					res.WriteHeader(200)
					fmt.Fprintf(res, "%s", `{"limit": 50, "first": {"href": "Href"}, "next": {"href": "Href"}, "virtual_private_endpoint_gateways": [{"href": "Href", "id": "ID", "name": "Name", "target_crn": "TargetCrn", "status": "Status"}]}`)
				}))
			})
			It(`Invoke ListVirtualPrivateEndpointGateways successfully`, func() {
				mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1(&mqcloudv1.MqcloudV1Options{
					URL:            testServer.URL,
					Authenticator:  &core.NoAuthAuthenticator{},
					AcceptLanguage: core.StringPtr(acceptLanguage),
				})
				Expect(serviceErr).To(BeNil())
				Expect(mqcloudService).ToNot(BeNil())

				// Invoke operation with nil options model (negative test)
				result, response, operationErr := mqcloudService.ListVirtualPrivateEndpointGateways(nil)
				Expect(operationErr).NotTo(BeNil())
				Expect(response).To(BeNil())
				Expect(result).To(BeNil())

				// Construct an instance of the ListVirtualPrivateEndpointGatewaysOptions model
				listVirtualPrivateEndpointGatewaysOptionsModel := new(mqcloudv1.ListVirtualPrivateEndpointGatewaysOptions)
				listVirtualPrivateEndpointGatewaysOptionsModel.ServiceInstanceGuid = core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				listVirtualPrivateEndpointGatewaysOptionsModel.TrustedProfile = core.StringPtr("crn:v1:bluemix:public:iam-identity::a/5d5ff2a9001c4055ab1408e9bf185f48::profile:Profile-1c0a8609-ca25-4ad2-a09b-aea472d34afc")
				listVirtualPrivateEndpointGatewaysOptionsModel.Start = core.StringPtr("r010-ebab3c08-c9a8-40c4-8869-61c09ddf7b44")
				listVirtualPrivateEndpointGatewaysOptionsModel.Limit = core.Int64Ptr(int64(10))
				listVirtualPrivateEndpointGatewaysOptionsModel.Headers = map[string]string{"x-custom-header": "x-custom-value"}

				// Invoke operation with valid options model (positive test)
				result, response, operationErr = mqcloudService.ListVirtualPrivateEndpointGateways(listVirtualPrivateEndpointGatewaysOptionsModel)
				Expect(operationErr).To(BeNil())
				Expect(response).ToNot(BeNil())
				Expect(result).ToNot(BeNil())

			})
			It(`Invoke ListVirtualPrivateEndpointGateways with error: Operation validation and request error`, func() {
				mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1(&mqcloudv1.MqcloudV1Options{
					URL:            testServer.URL,
					Authenticator:  &core.NoAuthAuthenticator{},
					AcceptLanguage: core.StringPtr(acceptLanguage),
				})
				Expect(serviceErr).To(BeNil())
				Expect(mqcloudService).ToNot(BeNil())

				// Construct an instance of the ListVirtualPrivateEndpointGatewaysOptions model
				listVirtualPrivateEndpointGatewaysOptionsModel := new(mqcloudv1.ListVirtualPrivateEndpointGatewaysOptions)
				listVirtualPrivateEndpointGatewaysOptionsModel.ServiceInstanceGuid = core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				listVirtualPrivateEndpointGatewaysOptionsModel.TrustedProfile = core.StringPtr("crn:v1:bluemix:public:iam-identity::a/5d5ff2a9001c4055ab1408e9bf185f48::profile:Profile-1c0a8609-ca25-4ad2-a09b-aea472d34afc")
				listVirtualPrivateEndpointGatewaysOptionsModel.Start = core.StringPtr("r010-ebab3c08-c9a8-40c4-8869-61c09ddf7b44")
				listVirtualPrivateEndpointGatewaysOptionsModel.Limit = core.Int64Ptr(int64(10))
				listVirtualPrivateEndpointGatewaysOptionsModel.Headers = map[string]string{"x-custom-header": "x-custom-value"}
				// Invoke operation with empty URL (negative test)
				err := mqcloudService.SetServiceURL("")
				Expect(err).To(BeNil())
				result, response, operationErr := mqcloudService.ListVirtualPrivateEndpointGateways(listVirtualPrivateEndpointGatewaysOptionsModel)
				Expect(operationErr).ToNot(BeNil())
				Expect(operationErr.Error()).To(ContainSubstring(core.ERRORMSG_SERVICE_URL_MISSING))
				Expect(response).To(BeNil())
				Expect(result).To(BeNil())
				// Construct a second instance of the ListVirtualPrivateEndpointGatewaysOptions model with no property values
				listVirtualPrivateEndpointGatewaysOptionsModelNew := new(mqcloudv1.ListVirtualPrivateEndpointGatewaysOptions)
				// Invoke operation with invalid model (negative test)
				result, response, operationErr = mqcloudService.ListVirtualPrivateEndpointGateways(listVirtualPrivateEndpointGatewaysOptionsModelNew)
				Expect(operationErr).ToNot(BeNil())
				Expect(response).To(BeNil())
				Expect(result).To(BeNil())
			})
			AfterEach(func() {
				testServer.Close()
			})
		})
		Context(`Using mock server endpoint with missing response body`, func() {
			BeforeEach(func() {
				testServer = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
					defer GinkgoRecover()

					// Set success status code with no respoonse body
					res.WriteHeader(200)
				}))
			})
			It(`Invoke ListVirtualPrivateEndpointGateways successfully`, func() {
				mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1(&mqcloudv1.MqcloudV1Options{
					URL:            testServer.URL,
					Authenticator:  &core.NoAuthAuthenticator{},
					AcceptLanguage: core.StringPtr(acceptLanguage),
				})
				Expect(serviceErr).To(BeNil())
				Expect(mqcloudService).ToNot(BeNil())

				// Construct an instance of the ListVirtualPrivateEndpointGatewaysOptions model
				listVirtualPrivateEndpointGatewaysOptionsModel := new(mqcloudv1.ListVirtualPrivateEndpointGatewaysOptions)
				listVirtualPrivateEndpointGatewaysOptionsModel.ServiceInstanceGuid = core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				listVirtualPrivateEndpointGatewaysOptionsModel.TrustedProfile = core.StringPtr("crn:v1:bluemix:public:iam-identity::a/5d5ff2a9001c4055ab1408e9bf185f48::profile:Profile-1c0a8609-ca25-4ad2-a09b-aea472d34afc")
				listVirtualPrivateEndpointGatewaysOptionsModel.Start = core.StringPtr("r010-ebab3c08-c9a8-40c4-8869-61c09ddf7b44")
				listVirtualPrivateEndpointGatewaysOptionsModel.Limit = core.Int64Ptr(int64(10))
				listVirtualPrivateEndpointGatewaysOptionsModel.Headers = map[string]string{"x-custom-header": "x-custom-value"}

				// Invoke operation
				result, response, operationErr := mqcloudService.ListVirtualPrivateEndpointGateways(listVirtualPrivateEndpointGatewaysOptionsModel)
				Expect(operationErr).To(BeNil())
				Expect(response).ToNot(BeNil())

				// Verify a nil result
				Expect(result).To(BeNil())
			})
			AfterEach(func() {
				testServer.Close()
			})
		})
		Context(`Test pagination helper method on response`, func() {
			It(`Invoke GetNextStart successfully`, func() {
				responseObject := new(mqcloudv1.VirtualPrivateEndpointGatewayDetailsCollection)
				nextObject := new(mqcloudv1.Next)
				nextObject.Href = core.StringPtr("ibm.com?start=abc-123")
				responseObject.Next = nextObject

				value, err := responseObject.GetNextStart()
				Expect(err).To(BeNil())
				Expect(value).To(Equal(core.StringPtr("abc-123")))
			})
			It(`Invoke GetNextStart without a "Next" property in the response`, func() {
				responseObject := new(mqcloudv1.VirtualPrivateEndpointGatewayDetailsCollection)

				value, err := responseObject.GetNextStart()
				Expect(err).To(BeNil())
				Expect(value).To(BeNil())
			})
			It(`Invoke GetNextStart without any query params in the "Next" URL`, func() {
				responseObject := new(mqcloudv1.VirtualPrivateEndpointGatewayDetailsCollection)
				nextObject := new(mqcloudv1.Next)
				nextObject.Href = core.StringPtr("ibm.com")
				responseObject.Next = nextObject

				value, err := responseObject.GetNextStart()
				Expect(err).To(BeNil())
				Expect(value).To(BeNil())
			})
		})
		Context(`Using mock server endpoint - paginated response`, func() {
			BeforeEach(func() {
				var requestNumber int = 0
				testServer = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
					defer GinkgoRecover()

					// Verify the contents of the request
					Expect(req.URL.EscapedPath()).To(Equal(listVirtualPrivateEndpointGatewaysPath))
					Expect(req.Method).To(Equal("GET"))

					// Set mock response
					res.Header().Set("Content-type", "application/json")
					res.WriteHeader(200)
					requestNumber++
					if requestNumber == 1 {
						fmt.Fprintf(res, "%s", `{"next":{"href":"https://myhost.com/somePath?start=1"},"total_count":2,"limit":1,"virtual_private_endpoint_gateways":[{"href":"Href","id":"ID","name":"Name","target_crn":"TargetCrn","status":"Status"}]}`)
					} else if requestNumber == 2 {
						fmt.Fprintf(res, "%s", `{"total_count":2,"limit":1,"virtual_private_endpoint_gateways":[{"href":"Href","id":"ID","name":"Name","target_crn":"TargetCrn","status":"Status"}]}`)
					} else {
						res.WriteHeader(400)
					}
				}))
			})
			It(`Use VirtualPrivateEndpointGatewaysPager.GetNext successfully`, func() {
				mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1(&mqcloudv1.MqcloudV1Options{
					URL:            testServer.URL,
					Authenticator:  &core.NoAuthAuthenticator{},
					AcceptLanguage: core.StringPtr(acceptLanguage),
				})
				Expect(serviceErr).To(BeNil())
				Expect(mqcloudService).ToNot(BeNil())

				listVirtualPrivateEndpointGatewaysOptionsModel := &mqcloudv1.ListVirtualPrivateEndpointGatewaysOptions{
					ServiceInstanceGuid: core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8"),
					TrustedProfile:      core.StringPtr("crn:v1:bluemix:public:iam-identity::a/5d5ff2a9001c4055ab1408e9bf185f48::profile:Profile-1c0a8609-ca25-4ad2-a09b-aea472d34afc"),
					Limit:               core.Int64Ptr(int64(10)),
				}

				pager, err := mqcloudService.NewVirtualPrivateEndpointGatewaysPager(listVirtualPrivateEndpointGatewaysOptionsModel)
				Expect(err).To(BeNil())
				Expect(pager).ToNot(BeNil())

				var allResults []mqcloudv1.VirtualPrivateEndpointGatewayDetails
				for pager.HasNext() {
					nextPage, err := pager.GetNext()
					Expect(err).To(BeNil())
					Expect(nextPage).ToNot(BeNil())
					allResults = append(allResults, nextPage...)
				}
				Expect(len(allResults)).To(Equal(2))
			})
			It(`Use VirtualPrivateEndpointGatewaysPager.GetAll successfully`, func() {
				mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1(&mqcloudv1.MqcloudV1Options{
					URL:            testServer.URL,
					Authenticator:  &core.NoAuthAuthenticator{},
					AcceptLanguage: core.StringPtr(acceptLanguage),
				})
				Expect(serviceErr).To(BeNil())
				Expect(mqcloudService).ToNot(BeNil())

				listVirtualPrivateEndpointGatewaysOptionsModel := &mqcloudv1.ListVirtualPrivateEndpointGatewaysOptions{
					ServiceInstanceGuid: core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8"),
					TrustedProfile:      core.StringPtr("crn:v1:bluemix:public:iam-identity::a/5d5ff2a9001c4055ab1408e9bf185f48::profile:Profile-1c0a8609-ca25-4ad2-a09b-aea472d34afc"),
					Limit:               core.Int64Ptr(int64(10)),
				}

				pager, err := mqcloudService.NewVirtualPrivateEndpointGatewaysPager(listVirtualPrivateEndpointGatewaysOptionsModel)
				Expect(err).To(BeNil())
				Expect(pager).ToNot(BeNil())

				allResults, err := pager.GetAll()
				Expect(err).To(BeNil())
				Expect(allResults).ToNot(BeNil())
				Expect(len(allResults)).To(Equal(2))
			})
		})
	})
	Describe(`GetVirtualPrivateEndpointGateway(getVirtualPrivateEndpointGatewayOptions *GetVirtualPrivateEndpointGatewayOptions) - Operation response error`, func() {
		acceptLanguage := "en-US,en;q=0.5"
		getVirtualPrivateEndpointGatewayPath := "/v1/a2b4d4bc-dadb-4637-bcec-9b7d1e723af8/virtual_private_endpoint_gateway/r010-ebab3c08-c9a8-40c4-8869-61c09ddf7b44"
		Context(`Using mock server endpoint with invalid JSON response`, func() {
			BeforeEach(func() {
				testServer = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
					defer GinkgoRecover()

					// Verify the contents of the request
					Expect(req.URL.EscapedPath()).To(Equal(getVirtualPrivateEndpointGatewayPath))
					Expect(req.Method).To(Equal("GET"))
					Expect(req.Header["Accept-Language"]).ToNot(BeNil())
					Expect(req.Header["Accept-Language"][0]).To(Equal(fmt.Sprintf("%v", "en-US,en;q=0.5")))
					Expect(req.Header["Trusted-Profile"]).ToNot(BeNil())
					Expect(req.Header["Trusted-Profile"][0]).To(Equal(fmt.Sprintf("%v", "crn:v1:bluemix:public:iam-identity::a/5d5ff2a9001c4055ab1408e9bf185f48::profile:Profile-1c0a8609-ca25-4ad2-a09b-aea472d34afc")))
					res.Header().Set("Content-type", "application/json")
					res.WriteHeader(200)
					fmt.Fprint(res, `} this is not valid json {`)
				}))
			})
			It(`Invoke GetVirtualPrivateEndpointGateway with error: Operation response processing error`, func() {
				mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1(&mqcloudv1.MqcloudV1Options{
					URL:            testServer.URL,
					Authenticator:  &core.NoAuthAuthenticator{},
					AcceptLanguage: core.StringPtr(acceptLanguage),
				})
				Expect(serviceErr).To(BeNil())
				Expect(mqcloudService).ToNot(BeNil())

				// Construct an instance of the GetVirtualPrivateEndpointGatewayOptions model
				getVirtualPrivateEndpointGatewayOptionsModel := new(mqcloudv1.GetVirtualPrivateEndpointGatewayOptions)
				getVirtualPrivateEndpointGatewayOptionsModel.ServiceInstanceGuid = core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				getVirtualPrivateEndpointGatewayOptionsModel.VirtualPrivateEndpointGatewayGuid = core.StringPtr("r010-ebab3c08-c9a8-40c4-8869-61c09ddf7b44")
				getVirtualPrivateEndpointGatewayOptionsModel.TrustedProfile = core.StringPtr("crn:v1:bluemix:public:iam-identity::a/5d5ff2a9001c4055ab1408e9bf185f48::profile:Profile-1c0a8609-ca25-4ad2-a09b-aea472d34afc")
				getVirtualPrivateEndpointGatewayOptionsModel.Headers = map[string]string{"x-custom-header": "x-custom-value"}
				// Expect response parsing to fail since we are receiving a text/plain response
				result, response, operationErr := mqcloudService.GetVirtualPrivateEndpointGateway(getVirtualPrivateEndpointGatewayOptionsModel)
				Expect(operationErr).ToNot(BeNil())
				Expect(response).ToNot(BeNil())
				Expect(result).To(BeNil())

				// Enable retries and test again
				mqcloudService.EnableRetries(0, 0)
				result, response, operationErr = mqcloudService.GetVirtualPrivateEndpointGateway(getVirtualPrivateEndpointGatewayOptionsModel)
				Expect(operationErr).ToNot(BeNil())
				Expect(response).ToNot(BeNil())
				Expect(result).To(BeNil())
			})
			AfterEach(func() {
				testServer.Close()
			})
		})
	})
	Describe(`GetVirtualPrivateEndpointGateway(getVirtualPrivateEndpointGatewayOptions *GetVirtualPrivateEndpointGatewayOptions)`, func() {
		acceptLanguage := "en-US,en;q=0.5"
		getVirtualPrivateEndpointGatewayPath := "/v1/a2b4d4bc-dadb-4637-bcec-9b7d1e723af8/virtual_private_endpoint_gateway/r010-ebab3c08-c9a8-40c4-8869-61c09ddf7b44"
		Context(`Using mock server endpoint with timeout`, func() {
			BeforeEach(func() {
				testServer = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
					defer GinkgoRecover()

					// Verify the contents of the request
					Expect(req.URL.EscapedPath()).To(Equal(getVirtualPrivateEndpointGatewayPath))
					Expect(req.Method).To(Equal("GET"))

					Expect(req.Header["Accept-Language"]).ToNot(BeNil())
					Expect(req.Header["Accept-Language"][0]).To(Equal(fmt.Sprintf("%v", "en-US,en;q=0.5")))
					Expect(req.Header["Trusted-Profile"]).ToNot(BeNil())
					Expect(req.Header["Trusted-Profile"][0]).To(Equal(fmt.Sprintf("%v", "crn:v1:bluemix:public:iam-identity::a/5d5ff2a9001c4055ab1408e9bf185f48::profile:Profile-1c0a8609-ca25-4ad2-a09b-aea472d34afc")))
					// Sleep a short time to support a timeout test
					time.Sleep(100 * time.Millisecond)

					// Set mock response
					res.Header().Set("Content-type", "application/json")
					res.WriteHeader(200)
					fmt.Fprintf(res, "%s", `{"href": "Href", "id": "ID", "name": "Name", "target_crn": "TargetCrn", "status": "Status"}`)
				}))
			})
			It(`Invoke GetVirtualPrivateEndpointGateway successfully with retries`, func() {
				mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1(&mqcloudv1.MqcloudV1Options{
					URL:            testServer.URL,
					Authenticator:  &core.NoAuthAuthenticator{},
					AcceptLanguage: core.StringPtr(acceptLanguage),
				})
				Expect(serviceErr).To(BeNil())
				Expect(mqcloudService).ToNot(BeNil())
				mqcloudService.EnableRetries(0, 0)

				// Construct an instance of the GetVirtualPrivateEndpointGatewayOptions model
				getVirtualPrivateEndpointGatewayOptionsModel := new(mqcloudv1.GetVirtualPrivateEndpointGatewayOptions)
				getVirtualPrivateEndpointGatewayOptionsModel.ServiceInstanceGuid = core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				getVirtualPrivateEndpointGatewayOptionsModel.VirtualPrivateEndpointGatewayGuid = core.StringPtr("r010-ebab3c08-c9a8-40c4-8869-61c09ddf7b44")
				getVirtualPrivateEndpointGatewayOptionsModel.TrustedProfile = core.StringPtr("crn:v1:bluemix:public:iam-identity::a/5d5ff2a9001c4055ab1408e9bf185f48::profile:Profile-1c0a8609-ca25-4ad2-a09b-aea472d34afc")
				getVirtualPrivateEndpointGatewayOptionsModel.Headers = map[string]string{"x-custom-header": "x-custom-value"}

				// Invoke operation with a Context to test a timeout error
				ctx, cancelFunc := context.WithTimeout(context.Background(), 80*time.Millisecond)
				defer cancelFunc()
				_, _, operationErr := mqcloudService.GetVirtualPrivateEndpointGatewayWithContext(ctx, getVirtualPrivateEndpointGatewayOptionsModel)
				Expect(operationErr).ToNot(BeNil())
				Expect(operationErr.Error()).To(ContainSubstring("deadline exceeded"))

				// Disable retries and test again
				mqcloudService.DisableRetries()
				result, response, operationErr := mqcloudService.GetVirtualPrivateEndpointGateway(getVirtualPrivateEndpointGatewayOptionsModel)
				Expect(operationErr).To(BeNil())
				Expect(response).ToNot(BeNil())
				Expect(result).ToNot(BeNil())

				// Re-test the timeout error with retries disabled
				ctx, cancelFunc2 := context.WithTimeout(context.Background(), 80*time.Millisecond)
				defer cancelFunc2()
				_, _, operationErr = mqcloudService.GetVirtualPrivateEndpointGatewayWithContext(ctx, getVirtualPrivateEndpointGatewayOptionsModel)
				Expect(operationErr).ToNot(BeNil())
				Expect(operationErr.Error()).To(ContainSubstring("deadline exceeded"))
			})
			AfterEach(func() {
				testServer.Close()
			})
		})
		Context(`Using mock server endpoint`, func() {
			BeforeEach(func() {
				testServer = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
					defer GinkgoRecover()

					// Verify the contents of the request
					Expect(req.URL.EscapedPath()).To(Equal(getVirtualPrivateEndpointGatewayPath))
					Expect(req.Method).To(Equal("GET"))

					Expect(req.Header["Accept-Language"]).ToNot(BeNil())
					Expect(req.Header["Accept-Language"][0]).To(Equal(fmt.Sprintf("%v", "en-US,en;q=0.5")))
					Expect(req.Header["Trusted-Profile"]).ToNot(BeNil())
					Expect(req.Header["Trusted-Profile"][0]).To(Equal(fmt.Sprintf("%v", "crn:v1:bluemix:public:iam-identity::a/5d5ff2a9001c4055ab1408e9bf185f48::profile:Profile-1c0a8609-ca25-4ad2-a09b-aea472d34afc")))
					// Set mock response
					res.Header().Set("Content-type", "application/json")
					res.WriteHeader(200)
					fmt.Fprintf(res, "%s", `{"href": "Href", "id": "ID", "name": "Name", "target_crn": "TargetCrn", "status": "Status"}`)
				}))
			})
			It(`Invoke GetVirtualPrivateEndpointGateway successfully`, func() {
				mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1(&mqcloudv1.MqcloudV1Options{
					URL:            testServer.URL,
					Authenticator:  &core.NoAuthAuthenticator{},
					AcceptLanguage: core.StringPtr(acceptLanguage),
				})
				Expect(serviceErr).To(BeNil())
				Expect(mqcloudService).ToNot(BeNil())

				// Invoke operation with nil options model (negative test)
				result, response, operationErr := mqcloudService.GetVirtualPrivateEndpointGateway(nil)
				Expect(operationErr).NotTo(BeNil())
				Expect(response).To(BeNil())
				Expect(result).To(BeNil())

				// Construct an instance of the GetVirtualPrivateEndpointGatewayOptions model
				getVirtualPrivateEndpointGatewayOptionsModel := new(mqcloudv1.GetVirtualPrivateEndpointGatewayOptions)
				getVirtualPrivateEndpointGatewayOptionsModel.ServiceInstanceGuid = core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				getVirtualPrivateEndpointGatewayOptionsModel.VirtualPrivateEndpointGatewayGuid = core.StringPtr("r010-ebab3c08-c9a8-40c4-8869-61c09ddf7b44")
				getVirtualPrivateEndpointGatewayOptionsModel.TrustedProfile = core.StringPtr("crn:v1:bluemix:public:iam-identity::a/5d5ff2a9001c4055ab1408e9bf185f48::profile:Profile-1c0a8609-ca25-4ad2-a09b-aea472d34afc")
				getVirtualPrivateEndpointGatewayOptionsModel.Headers = map[string]string{"x-custom-header": "x-custom-value"}

				// Invoke operation with valid options model (positive test)
				result, response, operationErr = mqcloudService.GetVirtualPrivateEndpointGateway(getVirtualPrivateEndpointGatewayOptionsModel)
				Expect(operationErr).To(BeNil())
				Expect(response).ToNot(BeNil())
				Expect(result).ToNot(BeNil())

			})
			It(`Invoke GetVirtualPrivateEndpointGateway with error: Operation validation and request error`, func() {
				mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1(&mqcloudv1.MqcloudV1Options{
					URL:            testServer.URL,
					Authenticator:  &core.NoAuthAuthenticator{},
					AcceptLanguage: core.StringPtr(acceptLanguage),
				})
				Expect(serviceErr).To(BeNil())
				Expect(mqcloudService).ToNot(BeNil())

				// Construct an instance of the GetVirtualPrivateEndpointGatewayOptions model
				getVirtualPrivateEndpointGatewayOptionsModel := new(mqcloudv1.GetVirtualPrivateEndpointGatewayOptions)
				getVirtualPrivateEndpointGatewayOptionsModel.ServiceInstanceGuid = core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				getVirtualPrivateEndpointGatewayOptionsModel.VirtualPrivateEndpointGatewayGuid = core.StringPtr("r010-ebab3c08-c9a8-40c4-8869-61c09ddf7b44")
				getVirtualPrivateEndpointGatewayOptionsModel.TrustedProfile = core.StringPtr("crn:v1:bluemix:public:iam-identity::a/5d5ff2a9001c4055ab1408e9bf185f48::profile:Profile-1c0a8609-ca25-4ad2-a09b-aea472d34afc")
				getVirtualPrivateEndpointGatewayOptionsModel.Headers = map[string]string{"x-custom-header": "x-custom-value"}
				// Invoke operation with empty URL (negative test)
				err := mqcloudService.SetServiceURL("")
				Expect(err).To(BeNil())
				result, response, operationErr := mqcloudService.GetVirtualPrivateEndpointGateway(getVirtualPrivateEndpointGatewayOptionsModel)
				Expect(operationErr).ToNot(BeNil())
				Expect(operationErr.Error()).To(ContainSubstring(core.ERRORMSG_SERVICE_URL_MISSING))
				Expect(response).To(BeNil())
				Expect(result).To(BeNil())
				// Construct a second instance of the GetVirtualPrivateEndpointGatewayOptions model with no property values
				getVirtualPrivateEndpointGatewayOptionsModelNew := new(mqcloudv1.GetVirtualPrivateEndpointGatewayOptions)
				// Invoke operation with invalid model (negative test)
				result, response, operationErr = mqcloudService.GetVirtualPrivateEndpointGateway(getVirtualPrivateEndpointGatewayOptionsModelNew)
				Expect(operationErr).ToNot(BeNil())
				Expect(response).To(BeNil())
				Expect(result).To(BeNil())
			})
			AfterEach(func() {
				testServer.Close()
			})
		})
		Context(`Using mock server endpoint with missing response body`, func() {
			BeforeEach(func() {
				testServer = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
					defer GinkgoRecover()

					// Set success status code with no respoonse body
					res.WriteHeader(200)
				}))
			})
			It(`Invoke GetVirtualPrivateEndpointGateway successfully`, func() {
				mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1(&mqcloudv1.MqcloudV1Options{
					URL:            testServer.URL,
					Authenticator:  &core.NoAuthAuthenticator{},
					AcceptLanguage: core.StringPtr(acceptLanguage),
				})
				Expect(serviceErr).To(BeNil())
				Expect(mqcloudService).ToNot(BeNil())

				// Construct an instance of the GetVirtualPrivateEndpointGatewayOptions model
				getVirtualPrivateEndpointGatewayOptionsModel := new(mqcloudv1.GetVirtualPrivateEndpointGatewayOptions)
				getVirtualPrivateEndpointGatewayOptionsModel.ServiceInstanceGuid = core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				getVirtualPrivateEndpointGatewayOptionsModel.VirtualPrivateEndpointGatewayGuid = core.StringPtr("r010-ebab3c08-c9a8-40c4-8869-61c09ddf7b44")
				getVirtualPrivateEndpointGatewayOptionsModel.TrustedProfile = core.StringPtr("crn:v1:bluemix:public:iam-identity::a/5d5ff2a9001c4055ab1408e9bf185f48::profile:Profile-1c0a8609-ca25-4ad2-a09b-aea472d34afc")
				getVirtualPrivateEndpointGatewayOptionsModel.Headers = map[string]string{"x-custom-header": "x-custom-value"}

				// Invoke operation
				result, response, operationErr := mqcloudService.GetVirtualPrivateEndpointGateway(getVirtualPrivateEndpointGatewayOptionsModel)
				Expect(operationErr).To(BeNil())
				Expect(response).ToNot(BeNil())

				// Verify a nil result
				Expect(result).To(BeNil())
			})
			AfterEach(func() {
				testServer.Close()
			})
		})
	})
	Describe(`DeleteVirtualPrivateEndpointGateway(deleteVirtualPrivateEndpointGatewayOptions *DeleteVirtualPrivateEndpointGatewayOptions)`, func() {
		acceptLanguage := "en-US,en;q=0.5"
		deleteVirtualPrivateEndpointGatewayPath := "/v1/a2b4d4bc-dadb-4637-bcec-9b7d1e723af8/virtual_private_endpoint_gateway/r010-ebab3c08-c9a8-40c4-8869-61c09ddf7b44"
		Context(`Using mock server endpoint`, func() {
			BeforeEach(func() {
				testServer = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
					defer GinkgoRecover()

					// Verify the contents of the request
					Expect(req.URL.EscapedPath()).To(Equal(deleteVirtualPrivateEndpointGatewayPath))
					Expect(req.Method).To(Equal("DELETE"))

					Expect(req.Header["Accept-Language"]).ToNot(BeNil())
					Expect(req.Header["Accept-Language"][0]).To(Equal(fmt.Sprintf("%v", "en-US,en;q=0.5")))
					Expect(req.Header["Trusted-Profile"]).ToNot(BeNil())
					Expect(req.Header["Trusted-Profile"][0]).To(Equal(fmt.Sprintf("%v", "crn:v1:bluemix:public:iam-identity::a/5d5ff2a9001c4055ab1408e9bf185f48::profile:Profile-1c0a8609-ca25-4ad2-a09b-aea472d34afc")))
					res.WriteHeader(204)
				}))
			})
			It(`Invoke DeleteVirtualPrivateEndpointGateway successfully`, func() {
				mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1(&mqcloudv1.MqcloudV1Options{
					URL:            testServer.URL,
					Authenticator:  &core.NoAuthAuthenticator{},
					AcceptLanguage: core.StringPtr(acceptLanguage),
				})
				Expect(serviceErr).To(BeNil())
				Expect(mqcloudService).ToNot(BeNil())

				// Invoke operation with nil options model (negative test)
				response, operationErr := mqcloudService.DeleteVirtualPrivateEndpointGateway(nil)
				Expect(operationErr).NotTo(BeNil())
				Expect(response).To(BeNil())

				// Construct an instance of the DeleteVirtualPrivateEndpointGatewayOptions model
				deleteVirtualPrivateEndpointGatewayOptionsModel := new(mqcloudv1.DeleteVirtualPrivateEndpointGatewayOptions)
				deleteVirtualPrivateEndpointGatewayOptionsModel.ServiceInstanceGuid = core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				deleteVirtualPrivateEndpointGatewayOptionsModel.VirtualPrivateEndpointGatewayGuid = core.StringPtr("r010-ebab3c08-c9a8-40c4-8869-61c09ddf7b44")
				deleteVirtualPrivateEndpointGatewayOptionsModel.TrustedProfile = core.StringPtr("crn:v1:bluemix:public:iam-identity::a/5d5ff2a9001c4055ab1408e9bf185f48::profile:Profile-1c0a8609-ca25-4ad2-a09b-aea472d34afc")
				deleteVirtualPrivateEndpointGatewayOptionsModel.Headers = map[string]string{"x-custom-header": "x-custom-value"}

				// Invoke operation with valid options model (positive test)
				response, operationErr = mqcloudService.DeleteVirtualPrivateEndpointGateway(deleteVirtualPrivateEndpointGatewayOptionsModel)
				Expect(operationErr).To(BeNil())
				Expect(response).ToNot(BeNil())
			})
			It(`Invoke DeleteVirtualPrivateEndpointGateway with error: Operation validation and request error`, func() {
				mqcloudService, serviceErr := mqcloudv1.NewMqcloudV1(&mqcloudv1.MqcloudV1Options{
					URL:            testServer.URL,
					Authenticator:  &core.NoAuthAuthenticator{},
					AcceptLanguage: core.StringPtr(acceptLanguage),
				})
				Expect(serviceErr).To(BeNil())
				Expect(mqcloudService).ToNot(BeNil())

				// Construct an instance of the DeleteVirtualPrivateEndpointGatewayOptions model
				deleteVirtualPrivateEndpointGatewayOptionsModel := new(mqcloudv1.DeleteVirtualPrivateEndpointGatewayOptions)
				deleteVirtualPrivateEndpointGatewayOptionsModel.ServiceInstanceGuid = core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				deleteVirtualPrivateEndpointGatewayOptionsModel.VirtualPrivateEndpointGatewayGuid = core.StringPtr("r010-ebab3c08-c9a8-40c4-8869-61c09ddf7b44")
				deleteVirtualPrivateEndpointGatewayOptionsModel.TrustedProfile = core.StringPtr("crn:v1:bluemix:public:iam-identity::a/5d5ff2a9001c4055ab1408e9bf185f48::profile:Profile-1c0a8609-ca25-4ad2-a09b-aea472d34afc")
				deleteVirtualPrivateEndpointGatewayOptionsModel.Headers = map[string]string{"x-custom-header": "x-custom-value"}
				// Invoke operation with empty URL (negative test)
				err := mqcloudService.SetServiceURL("")
				Expect(err).To(BeNil())
				response, operationErr := mqcloudService.DeleteVirtualPrivateEndpointGateway(deleteVirtualPrivateEndpointGatewayOptionsModel)
				Expect(operationErr).ToNot(BeNil())
				Expect(operationErr.Error()).To(ContainSubstring(core.ERRORMSG_SERVICE_URL_MISSING))
				Expect(response).To(BeNil())
				// Construct a second instance of the DeleteVirtualPrivateEndpointGatewayOptions model with no property values
				deleteVirtualPrivateEndpointGatewayOptionsModelNew := new(mqcloudv1.DeleteVirtualPrivateEndpointGatewayOptions)
				// Invoke operation with invalid model (negative test)
				response, operationErr = mqcloudService.DeleteVirtualPrivateEndpointGateway(deleteVirtualPrivateEndpointGatewayOptionsModelNew)
				Expect(operationErr).ToNot(BeNil())
				Expect(response).To(BeNil())
			})
			AfterEach(func() {
				testServer.Close()
			})
		})
	})
	Describe(`Model constructor tests`, func() {
		Context(`Using a service client instance`, func() {
			acceptLanguage := "en-US,en;q=0.5"
			mqcloudService, _ := mqcloudv1.NewMqcloudV1(&mqcloudv1.MqcloudV1Options{
				URL:            "http://mqcloudv1modelgenerator.com",
				Authenticator:  &core.NoAuthAuthenticator{},
				AcceptLanguage: core.StringPtr(acceptLanguage),
			})
			It(`Invoke NewCreateApplicationApikeyOptions successfully`, func() {
				// Construct an instance of the CreateApplicationApikeyOptions model
				serviceInstanceGuid := "a2b4d4bc-dadb-4637-bcec-9b7d1e723af8"
				applicationID := "0123456789ABCDEF0123456789ABCDEF"
				createApplicationApikeyOptionsName := "test-api-key"
				createApplicationApikeyOptionsModel := mqcloudService.NewCreateApplicationApikeyOptions(serviceInstanceGuid, applicationID, createApplicationApikeyOptionsName)
				createApplicationApikeyOptionsModel.SetServiceInstanceGuid("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				createApplicationApikeyOptionsModel.SetApplicationID("0123456789ABCDEF0123456789ABCDEF")
				createApplicationApikeyOptionsModel.SetName("test-api-key")
				createApplicationApikeyOptionsModel.SetHeaders(map[string]string{"foo": "bar"})
				Expect(createApplicationApikeyOptionsModel).ToNot(BeNil())
				Expect(createApplicationApikeyOptionsModel.ServiceInstanceGuid).To(Equal(core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")))
				Expect(createApplicationApikeyOptionsModel.ApplicationID).To(Equal(core.StringPtr("0123456789ABCDEF0123456789ABCDEF")))
				Expect(createApplicationApikeyOptionsModel.Name).To(Equal(core.StringPtr("test-api-key")))
				Expect(createApplicationApikeyOptionsModel.Headers).To(Equal(map[string]string{"foo": "bar"}))
			})
			It(`Invoke NewCreateApplicationOptions successfully`, func() {
				// Construct an instance of the CreateApplicationOptions model
				serviceInstanceGuid := "a2b4d4bc-dadb-4637-bcec-9b7d1e723af8"
				createApplicationOptionsName := "test-app"
				createApplicationOptionsModel := mqcloudService.NewCreateApplicationOptions(serviceInstanceGuid, createApplicationOptionsName)
				createApplicationOptionsModel.SetServiceInstanceGuid("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				createApplicationOptionsModel.SetName("test-app")
				createApplicationOptionsModel.SetHeaders(map[string]string{"foo": "bar"})
				Expect(createApplicationOptionsModel).ToNot(BeNil())
				Expect(createApplicationOptionsModel.ServiceInstanceGuid).To(Equal(core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")))
				Expect(createApplicationOptionsModel.Name).To(Equal(core.StringPtr("test-app")))
				Expect(createApplicationOptionsModel.Headers).To(Equal(map[string]string{"foo": "bar"}))
			})
			It(`Invoke NewCreateKeyStorePemCertificateOptions successfully`, func() {
				// Construct an instance of the CreateKeyStorePemCertificateOptions model
				serviceInstanceGuid := "a2b4d4bc-dadb-4637-bcec-9b7d1e723af8"
				queueManagerID := "b8e1aeda078009cf3db74e90d5d42328"
				label := "certlabel"
				certificateFile := CreateMockReader("This is a mock file.")
				createKeyStorePemCertificateOptionsModel := mqcloudService.NewCreateKeyStorePemCertificateOptions(serviceInstanceGuid, queueManagerID, label, certificateFile)
				createKeyStorePemCertificateOptionsModel.SetServiceInstanceGuid("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				createKeyStorePemCertificateOptionsModel.SetQueueManagerID("b8e1aeda078009cf3db74e90d5d42328")
				createKeyStorePemCertificateOptionsModel.SetLabel("certlabel")
				createKeyStorePemCertificateOptionsModel.SetCertificateFile(CreateMockReader("This is a mock file."))
				createKeyStorePemCertificateOptionsModel.SetHeaders(map[string]string{"foo": "bar"})
				Expect(createKeyStorePemCertificateOptionsModel).ToNot(BeNil())
				Expect(createKeyStorePemCertificateOptionsModel.ServiceInstanceGuid).To(Equal(core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")))
				Expect(createKeyStorePemCertificateOptionsModel.QueueManagerID).To(Equal(core.StringPtr("b8e1aeda078009cf3db74e90d5d42328")))
				Expect(createKeyStorePemCertificateOptionsModel.Label).To(Equal(core.StringPtr("certlabel")))
				Expect(createKeyStorePemCertificateOptionsModel.CertificateFile).To(Equal(CreateMockReader("This is a mock file.")))
				Expect(createKeyStorePemCertificateOptionsModel.Headers).To(Equal(map[string]string{"foo": "bar"}))
			})
			It(`Invoke NewCreateQueueManagerOptions successfully`, func() {
				// Construct an instance of the CreateQueueManagerOptions model
				serviceInstanceGuid := "a2b4d4bc-dadb-4637-bcec-9b7d1e723af8"
				createQueueManagerOptionsName := "testqm"
				createQueueManagerOptionsLocation := "reserved-eu-de-cluster-f884"
				createQueueManagerOptionsSize := "small"
				createQueueManagerOptionsModel := mqcloudService.NewCreateQueueManagerOptions(serviceInstanceGuid, createQueueManagerOptionsName, createQueueManagerOptionsLocation, createQueueManagerOptionsSize)
				createQueueManagerOptionsModel.SetServiceInstanceGuid("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				createQueueManagerOptionsModel.SetName("testqm")
				createQueueManagerOptionsModel.SetLocation("reserved-eu-de-cluster-f884")
				createQueueManagerOptionsModel.SetSize("small")
				createQueueManagerOptionsModel.SetDisplayName("A test queue manager")
				createQueueManagerOptionsModel.SetVersion("9.3.2_2")
				createQueueManagerOptionsModel.SetHeaders(map[string]string{"foo": "bar"})
				Expect(createQueueManagerOptionsModel).ToNot(BeNil())
				Expect(createQueueManagerOptionsModel.ServiceInstanceGuid).To(Equal(core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")))
				Expect(createQueueManagerOptionsModel.Name).To(Equal(core.StringPtr("testqm")))
				Expect(createQueueManagerOptionsModel.Location).To(Equal(core.StringPtr("reserved-eu-de-cluster-f884")))
				Expect(createQueueManagerOptionsModel.Size).To(Equal(core.StringPtr("small")))
				Expect(createQueueManagerOptionsModel.DisplayName).To(Equal(core.StringPtr("A test queue manager")))
				Expect(createQueueManagerOptionsModel.Version).To(Equal(core.StringPtr("9.3.2_2")))
				Expect(createQueueManagerOptionsModel.Headers).To(Equal(map[string]string{"foo": "bar"}))
			})
			It(`Invoke NewCreateTrustStorePemCertificateOptions successfully`, func() {
				// Construct an instance of the CreateTrustStorePemCertificateOptions model
				serviceInstanceGuid := "a2b4d4bc-dadb-4637-bcec-9b7d1e723af8"
				queueManagerID := "b8e1aeda078009cf3db74e90d5d42328"
				label := "certlabel"
				certificateFile := CreateMockReader("This is a mock file.")
				createTrustStorePemCertificateOptionsModel := mqcloudService.NewCreateTrustStorePemCertificateOptions(serviceInstanceGuid, queueManagerID, label, certificateFile)
				createTrustStorePemCertificateOptionsModel.SetServiceInstanceGuid("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				createTrustStorePemCertificateOptionsModel.SetQueueManagerID("b8e1aeda078009cf3db74e90d5d42328")
				createTrustStorePemCertificateOptionsModel.SetLabel("certlabel")
				createTrustStorePemCertificateOptionsModel.SetCertificateFile(CreateMockReader("This is a mock file."))
				createTrustStorePemCertificateOptionsModel.SetHeaders(map[string]string{"foo": "bar"})
				Expect(createTrustStorePemCertificateOptionsModel).ToNot(BeNil())
				Expect(createTrustStorePemCertificateOptionsModel.ServiceInstanceGuid).To(Equal(core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")))
				Expect(createTrustStorePemCertificateOptionsModel.QueueManagerID).To(Equal(core.StringPtr("b8e1aeda078009cf3db74e90d5d42328")))
				Expect(createTrustStorePemCertificateOptionsModel.Label).To(Equal(core.StringPtr("certlabel")))
				Expect(createTrustStorePemCertificateOptionsModel.CertificateFile).To(Equal(CreateMockReader("This is a mock file.")))
				Expect(createTrustStorePemCertificateOptionsModel.Headers).To(Equal(map[string]string{"foo": "bar"}))
			})
			It(`Invoke NewCreateUserOptions successfully`, func() {
				// Construct an instance of the CreateUserOptions model
				serviceInstanceGuid := "a2b4d4bc-dadb-4637-bcec-9b7d1e723af8"
				createUserOptionsEmail := "testuser@ibm.com"
				createUserOptionsName := "testuser"
				createUserOptionsModel := mqcloudService.NewCreateUserOptions(serviceInstanceGuid, createUserOptionsEmail, createUserOptionsName)
				createUserOptionsModel.SetServiceInstanceGuid("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				createUserOptionsModel.SetEmail("testuser@ibm.com")
				createUserOptionsModel.SetName("testuser")
				createUserOptionsModel.SetHeaders(map[string]string{"foo": "bar"})
				Expect(createUserOptionsModel).ToNot(BeNil())
				Expect(createUserOptionsModel.ServiceInstanceGuid).To(Equal(core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")))
				Expect(createUserOptionsModel.Email).To(Equal(core.StringPtr("testuser@ibm.com")))
				Expect(createUserOptionsModel.Name).To(Equal(core.StringPtr("testuser")))
				Expect(createUserOptionsModel.Headers).To(Equal(map[string]string{"foo": "bar"}))
			})
			It(`Invoke NewCreateVirtualPrivateEndpointGatewayOptions successfully`, func() {
				// Construct an instance of the CreateVirtualPrivateEndpointGatewayOptions model
				serviceInstanceGuid := "a2b4d4bc-dadb-4637-bcec-9b7d1e723af8"
				createVirtualPrivateEndpointGatewayOptionsName := "vpe_gateway1-to-vpe_gateway2"
				createVirtualPrivateEndpointGatewayOptionsTargetCrn := "crn:v1:bluemix:public:mqcloud:eu-de:::endpoint:qm1.private.eu-de.mq2.test.appdomain.cloud"
				createVirtualPrivateEndpointGatewayOptionsModel := mqcloudService.NewCreateVirtualPrivateEndpointGatewayOptions(serviceInstanceGuid, createVirtualPrivateEndpointGatewayOptionsName, createVirtualPrivateEndpointGatewayOptionsTargetCrn)
				createVirtualPrivateEndpointGatewayOptionsModel.SetServiceInstanceGuid("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				createVirtualPrivateEndpointGatewayOptionsModel.SetName("vpe_gateway1-to-vpe_gateway2")
				createVirtualPrivateEndpointGatewayOptionsModel.SetTargetCrn("crn:v1:bluemix:public:mqcloud:eu-de:::endpoint:qm1.private.eu-de.mq2.test.appdomain.cloud")
				createVirtualPrivateEndpointGatewayOptionsModel.SetTrustedProfile("crn:v1:bluemix:public:iam-identity::a/5d5ff2a9001c4055ab1408e9bf185f48::profile:Profile-1c0a8609-ca25-4ad2-a09b-aea472d34afc")
				createVirtualPrivateEndpointGatewayOptionsModel.SetHeaders(map[string]string{"foo": "bar"})
				Expect(createVirtualPrivateEndpointGatewayOptionsModel).ToNot(BeNil())
				Expect(createVirtualPrivateEndpointGatewayOptionsModel.ServiceInstanceGuid).To(Equal(core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")))
				Expect(createVirtualPrivateEndpointGatewayOptionsModel.Name).To(Equal(core.StringPtr("vpe_gateway1-to-vpe_gateway2")))
				Expect(createVirtualPrivateEndpointGatewayOptionsModel.TargetCrn).To(Equal(core.StringPtr("crn:v1:bluemix:public:mqcloud:eu-de:::endpoint:qm1.private.eu-de.mq2.test.appdomain.cloud")))
				Expect(createVirtualPrivateEndpointGatewayOptionsModel.TrustedProfile).To(Equal(core.StringPtr("crn:v1:bluemix:public:iam-identity::a/5d5ff2a9001c4055ab1408e9bf185f48::profile:Profile-1c0a8609-ca25-4ad2-a09b-aea472d34afc")))
				Expect(createVirtualPrivateEndpointGatewayOptionsModel.Headers).To(Equal(map[string]string{"foo": "bar"}))
			})
			It(`Invoke NewDeleteApplicationOptions successfully`, func() {
				// Construct an instance of the DeleteApplicationOptions model
				serviceInstanceGuid := "a2b4d4bc-dadb-4637-bcec-9b7d1e723af8"
				applicationID := "0123456789ABCDEF0123456789ABCDEF"
				deleteApplicationOptionsModel := mqcloudService.NewDeleteApplicationOptions(serviceInstanceGuid, applicationID)
				deleteApplicationOptionsModel.SetServiceInstanceGuid("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				deleteApplicationOptionsModel.SetApplicationID("0123456789ABCDEF0123456789ABCDEF")
				deleteApplicationOptionsModel.SetHeaders(map[string]string{"foo": "bar"})
				Expect(deleteApplicationOptionsModel).ToNot(BeNil())
				Expect(deleteApplicationOptionsModel.ServiceInstanceGuid).To(Equal(core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")))
				Expect(deleteApplicationOptionsModel.ApplicationID).To(Equal(core.StringPtr("0123456789ABCDEF0123456789ABCDEF")))
				Expect(deleteApplicationOptionsModel.Headers).To(Equal(map[string]string{"foo": "bar"}))
			})
			It(`Invoke NewDeleteKeyStoreCertificateOptions successfully`, func() {
				// Construct an instance of the DeleteKeyStoreCertificateOptions model
				serviceInstanceGuid := "a2b4d4bc-dadb-4637-bcec-9b7d1e723af8"
				queueManagerID := "b8e1aeda078009cf3db74e90d5d42328"
				certificateID := "9b7d1e723af8233"
				deleteKeyStoreCertificateOptionsModel := mqcloudService.NewDeleteKeyStoreCertificateOptions(serviceInstanceGuid, queueManagerID, certificateID)
				deleteKeyStoreCertificateOptionsModel.SetServiceInstanceGuid("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				deleteKeyStoreCertificateOptionsModel.SetQueueManagerID("b8e1aeda078009cf3db74e90d5d42328")
				deleteKeyStoreCertificateOptionsModel.SetCertificateID("9b7d1e723af8233")
				deleteKeyStoreCertificateOptionsModel.SetHeaders(map[string]string{"foo": "bar"})
				Expect(deleteKeyStoreCertificateOptionsModel).ToNot(BeNil())
				Expect(deleteKeyStoreCertificateOptionsModel.ServiceInstanceGuid).To(Equal(core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")))
				Expect(deleteKeyStoreCertificateOptionsModel.QueueManagerID).To(Equal(core.StringPtr("b8e1aeda078009cf3db74e90d5d42328")))
				Expect(deleteKeyStoreCertificateOptionsModel.CertificateID).To(Equal(core.StringPtr("9b7d1e723af8233")))
				Expect(deleteKeyStoreCertificateOptionsModel.Headers).To(Equal(map[string]string{"foo": "bar"}))
			})
			It(`Invoke NewDeleteQueueManagerOptions successfully`, func() {
				// Construct an instance of the DeleteQueueManagerOptions model
				serviceInstanceGuid := "a2b4d4bc-dadb-4637-bcec-9b7d1e723af8"
				queueManagerID := "b8e1aeda078009cf3db74e90d5d42328"
				deleteQueueManagerOptionsModel := mqcloudService.NewDeleteQueueManagerOptions(serviceInstanceGuid, queueManagerID)
				deleteQueueManagerOptionsModel.SetServiceInstanceGuid("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				deleteQueueManagerOptionsModel.SetQueueManagerID("b8e1aeda078009cf3db74e90d5d42328")
				deleteQueueManagerOptionsModel.SetHeaders(map[string]string{"foo": "bar"})
				Expect(deleteQueueManagerOptionsModel).ToNot(BeNil())
				Expect(deleteQueueManagerOptionsModel.ServiceInstanceGuid).To(Equal(core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")))
				Expect(deleteQueueManagerOptionsModel.QueueManagerID).To(Equal(core.StringPtr("b8e1aeda078009cf3db74e90d5d42328")))
				Expect(deleteQueueManagerOptionsModel.Headers).To(Equal(map[string]string{"foo": "bar"}))
			})
			It(`Invoke NewDeleteTrustStoreCertificateOptions successfully`, func() {
				// Construct an instance of the DeleteTrustStoreCertificateOptions model
				serviceInstanceGuid := "a2b4d4bc-dadb-4637-bcec-9b7d1e723af8"
				queueManagerID := "b8e1aeda078009cf3db74e90d5d42328"
				certificateID := "9b7d1e723af8233"
				deleteTrustStoreCertificateOptionsModel := mqcloudService.NewDeleteTrustStoreCertificateOptions(serviceInstanceGuid, queueManagerID, certificateID)
				deleteTrustStoreCertificateOptionsModel.SetServiceInstanceGuid("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				deleteTrustStoreCertificateOptionsModel.SetQueueManagerID("b8e1aeda078009cf3db74e90d5d42328")
				deleteTrustStoreCertificateOptionsModel.SetCertificateID("9b7d1e723af8233")
				deleteTrustStoreCertificateOptionsModel.SetHeaders(map[string]string{"foo": "bar"})
				Expect(deleteTrustStoreCertificateOptionsModel).ToNot(BeNil())
				Expect(deleteTrustStoreCertificateOptionsModel.ServiceInstanceGuid).To(Equal(core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")))
				Expect(deleteTrustStoreCertificateOptionsModel.QueueManagerID).To(Equal(core.StringPtr("b8e1aeda078009cf3db74e90d5d42328")))
				Expect(deleteTrustStoreCertificateOptionsModel.CertificateID).To(Equal(core.StringPtr("9b7d1e723af8233")))
				Expect(deleteTrustStoreCertificateOptionsModel.Headers).To(Equal(map[string]string{"foo": "bar"}))
			})
			It(`Invoke NewDeleteUserOptions successfully`, func() {
				// Construct an instance of the DeleteUserOptions model
				serviceInstanceGuid := "a2b4d4bc-dadb-4637-bcec-9b7d1e723af8"
				userID := "31a413dd84346effc8895b6ba4641641"
				deleteUserOptionsModel := mqcloudService.NewDeleteUserOptions(serviceInstanceGuid, userID)
				deleteUserOptionsModel.SetServiceInstanceGuid("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				deleteUserOptionsModel.SetUserID("31a413dd84346effc8895b6ba4641641")
				deleteUserOptionsModel.SetHeaders(map[string]string{"foo": "bar"})
				Expect(deleteUserOptionsModel).ToNot(BeNil())
				Expect(deleteUserOptionsModel.ServiceInstanceGuid).To(Equal(core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")))
				Expect(deleteUserOptionsModel.UserID).To(Equal(core.StringPtr("31a413dd84346effc8895b6ba4641641")))
				Expect(deleteUserOptionsModel.Headers).To(Equal(map[string]string{"foo": "bar"}))
			})
			It(`Invoke NewDeleteVirtualPrivateEndpointGatewayOptions successfully`, func() {
				// Construct an instance of the DeleteVirtualPrivateEndpointGatewayOptions model
				serviceInstanceGuid := "a2b4d4bc-dadb-4637-bcec-9b7d1e723af8"
				virtualPrivateEndpointGatewayGuid := "r010-ebab3c08-c9a8-40c4-8869-61c09ddf7b44"
				deleteVirtualPrivateEndpointGatewayOptionsModel := mqcloudService.NewDeleteVirtualPrivateEndpointGatewayOptions(serviceInstanceGuid, virtualPrivateEndpointGatewayGuid)
				deleteVirtualPrivateEndpointGatewayOptionsModel.SetServiceInstanceGuid("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				deleteVirtualPrivateEndpointGatewayOptionsModel.SetVirtualPrivateEndpointGatewayGuid("r010-ebab3c08-c9a8-40c4-8869-61c09ddf7b44")
				deleteVirtualPrivateEndpointGatewayOptionsModel.SetTrustedProfile("crn:v1:bluemix:public:iam-identity::a/5d5ff2a9001c4055ab1408e9bf185f48::profile:Profile-1c0a8609-ca25-4ad2-a09b-aea472d34afc")
				deleteVirtualPrivateEndpointGatewayOptionsModel.SetHeaders(map[string]string{"foo": "bar"})
				Expect(deleteVirtualPrivateEndpointGatewayOptionsModel).ToNot(BeNil())
				Expect(deleteVirtualPrivateEndpointGatewayOptionsModel.ServiceInstanceGuid).To(Equal(core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")))
				Expect(deleteVirtualPrivateEndpointGatewayOptionsModel.VirtualPrivateEndpointGatewayGuid).To(Equal(core.StringPtr("r010-ebab3c08-c9a8-40c4-8869-61c09ddf7b44")))
				Expect(deleteVirtualPrivateEndpointGatewayOptionsModel.TrustedProfile).To(Equal(core.StringPtr("crn:v1:bluemix:public:iam-identity::a/5d5ff2a9001c4055ab1408e9bf185f48::profile:Profile-1c0a8609-ca25-4ad2-a09b-aea472d34afc")))
				Expect(deleteVirtualPrivateEndpointGatewayOptionsModel.Headers).To(Equal(map[string]string{"foo": "bar"}))
			})
			It(`Invoke NewDownloadKeyStoreCertificateOptions successfully`, func() {
				// Construct an instance of the DownloadKeyStoreCertificateOptions model
				serviceInstanceGuid := "a2b4d4bc-dadb-4637-bcec-9b7d1e723af8"
				queueManagerID := "b8e1aeda078009cf3db74e90d5d42328"
				certificateID := "9b7d1e723af8233"
				downloadKeyStoreCertificateOptionsModel := mqcloudService.NewDownloadKeyStoreCertificateOptions(serviceInstanceGuid, queueManagerID, certificateID)
				downloadKeyStoreCertificateOptionsModel.SetServiceInstanceGuid("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				downloadKeyStoreCertificateOptionsModel.SetQueueManagerID("b8e1aeda078009cf3db74e90d5d42328")
				downloadKeyStoreCertificateOptionsModel.SetCertificateID("9b7d1e723af8233")
				downloadKeyStoreCertificateOptionsModel.SetHeaders(map[string]string{"foo": "bar"})
				Expect(downloadKeyStoreCertificateOptionsModel).ToNot(BeNil())
				Expect(downloadKeyStoreCertificateOptionsModel.ServiceInstanceGuid).To(Equal(core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")))
				Expect(downloadKeyStoreCertificateOptionsModel.QueueManagerID).To(Equal(core.StringPtr("b8e1aeda078009cf3db74e90d5d42328")))
				Expect(downloadKeyStoreCertificateOptionsModel.CertificateID).To(Equal(core.StringPtr("9b7d1e723af8233")))
				Expect(downloadKeyStoreCertificateOptionsModel.Headers).To(Equal(map[string]string{"foo": "bar"}))
			})
			It(`Invoke NewDownloadTrustStoreCertificateOptions successfully`, func() {
				// Construct an instance of the DownloadTrustStoreCertificateOptions model
				serviceInstanceGuid := "a2b4d4bc-dadb-4637-bcec-9b7d1e723af8"
				queueManagerID := "b8e1aeda078009cf3db74e90d5d42328"
				certificateID := "9b7d1e723af8233"
				downloadTrustStoreCertificateOptionsModel := mqcloudService.NewDownloadTrustStoreCertificateOptions(serviceInstanceGuid, queueManagerID, certificateID)
				downloadTrustStoreCertificateOptionsModel.SetServiceInstanceGuid("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				downloadTrustStoreCertificateOptionsModel.SetQueueManagerID("b8e1aeda078009cf3db74e90d5d42328")
				downloadTrustStoreCertificateOptionsModel.SetCertificateID("9b7d1e723af8233")
				downloadTrustStoreCertificateOptionsModel.SetHeaders(map[string]string{"foo": "bar"})
				Expect(downloadTrustStoreCertificateOptionsModel).ToNot(BeNil())
				Expect(downloadTrustStoreCertificateOptionsModel.ServiceInstanceGuid).To(Equal(core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")))
				Expect(downloadTrustStoreCertificateOptionsModel.QueueManagerID).To(Equal(core.StringPtr("b8e1aeda078009cf3db74e90d5d42328")))
				Expect(downloadTrustStoreCertificateOptionsModel.CertificateID).To(Equal(core.StringPtr("9b7d1e723af8233")))
				Expect(downloadTrustStoreCertificateOptionsModel.Headers).To(Equal(map[string]string{"foo": "bar"}))
			})
			It(`Invoke NewGetApplicationOptions successfully`, func() {
				// Construct an instance of the GetApplicationOptions model
				serviceInstanceGuid := "a2b4d4bc-dadb-4637-bcec-9b7d1e723af8"
				applicationID := "0123456789ABCDEF0123456789ABCDEF"
				getApplicationOptionsModel := mqcloudService.NewGetApplicationOptions(serviceInstanceGuid, applicationID)
				getApplicationOptionsModel.SetServiceInstanceGuid("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				getApplicationOptionsModel.SetApplicationID("0123456789ABCDEF0123456789ABCDEF")
				getApplicationOptionsModel.SetHeaders(map[string]string{"foo": "bar"})
				Expect(getApplicationOptionsModel).ToNot(BeNil())
				Expect(getApplicationOptionsModel.ServiceInstanceGuid).To(Equal(core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")))
				Expect(getApplicationOptionsModel.ApplicationID).To(Equal(core.StringPtr("0123456789ABCDEF0123456789ABCDEF")))
				Expect(getApplicationOptionsModel.Headers).To(Equal(map[string]string{"foo": "bar"}))
			})
			It(`Invoke NewGetCertificateAmsChannelsOptions successfully`, func() {
				// Construct an instance of the GetCertificateAmsChannelsOptions model
				queueManagerID := "b8e1aeda078009cf3db74e90d5d42328"
				certificateID := "9b7d1e723af8233"
				serviceInstanceGuid := "a2b4d4bc-dadb-4637-bcec-9b7d1e723af8"
				getCertificateAmsChannelsOptionsModel := mqcloudService.NewGetCertificateAmsChannelsOptions(queueManagerID, certificateID, serviceInstanceGuid)
				getCertificateAmsChannelsOptionsModel.SetQueueManagerID("b8e1aeda078009cf3db74e90d5d42328")
				getCertificateAmsChannelsOptionsModel.SetCertificateID("9b7d1e723af8233")
				getCertificateAmsChannelsOptionsModel.SetServiceInstanceGuid("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				getCertificateAmsChannelsOptionsModel.SetHeaders(map[string]string{"foo": "bar"})
				Expect(getCertificateAmsChannelsOptionsModel).ToNot(BeNil())
				Expect(getCertificateAmsChannelsOptionsModel.QueueManagerID).To(Equal(core.StringPtr("b8e1aeda078009cf3db74e90d5d42328")))
				Expect(getCertificateAmsChannelsOptionsModel.CertificateID).To(Equal(core.StringPtr("9b7d1e723af8233")))
				Expect(getCertificateAmsChannelsOptionsModel.ServiceInstanceGuid).To(Equal(core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")))
				Expect(getCertificateAmsChannelsOptionsModel.Headers).To(Equal(map[string]string{"foo": "bar"}))
			})
			It(`Invoke NewGetKeyStoreCertificateOptions successfully`, func() {
				// Construct an instance of the GetKeyStoreCertificateOptions model
				serviceInstanceGuid := "a2b4d4bc-dadb-4637-bcec-9b7d1e723af8"
				queueManagerID := "b8e1aeda078009cf3db74e90d5d42328"
				certificateID := "9b7d1e723af8233"
				getKeyStoreCertificateOptionsModel := mqcloudService.NewGetKeyStoreCertificateOptions(serviceInstanceGuid, queueManagerID, certificateID)
				getKeyStoreCertificateOptionsModel.SetServiceInstanceGuid("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				getKeyStoreCertificateOptionsModel.SetQueueManagerID("b8e1aeda078009cf3db74e90d5d42328")
				getKeyStoreCertificateOptionsModel.SetCertificateID("9b7d1e723af8233")
				getKeyStoreCertificateOptionsModel.SetHeaders(map[string]string{"foo": "bar"})
				Expect(getKeyStoreCertificateOptionsModel).ToNot(BeNil())
				Expect(getKeyStoreCertificateOptionsModel.ServiceInstanceGuid).To(Equal(core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")))
				Expect(getKeyStoreCertificateOptionsModel.QueueManagerID).To(Equal(core.StringPtr("b8e1aeda078009cf3db74e90d5d42328")))
				Expect(getKeyStoreCertificateOptionsModel.CertificateID).To(Equal(core.StringPtr("9b7d1e723af8233")))
				Expect(getKeyStoreCertificateOptionsModel.Headers).To(Equal(map[string]string{"foo": "bar"}))
			})
			It(`Invoke NewGetOptionsOptions successfully`, func() {
				// Construct an instance of the GetOptionsOptions model
				serviceInstanceGuid := "a2b4d4bc-dadb-4637-bcec-9b7d1e723af8"
				getOptionsOptionsModel := mqcloudService.NewGetOptionsOptions(serviceInstanceGuid)
				getOptionsOptionsModel.SetServiceInstanceGuid("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				getOptionsOptionsModel.SetHeaders(map[string]string{"foo": "bar"})
				Expect(getOptionsOptionsModel).ToNot(BeNil())
				Expect(getOptionsOptionsModel.ServiceInstanceGuid).To(Equal(core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")))
				Expect(getOptionsOptionsModel.Headers).To(Equal(map[string]string{"foo": "bar"}))
			})
			It(`Invoke NewGetQueueManagerAvailableUpgradeVersionsOptions successfully`, func() {
				// Construct an instance of the GetQueueManagerAvailableUpgradeVersionsOptions model
				serviceInstanceGuid := "a2b4d4bc-dadb-4637-bcec-9b7d1e723af8"
				queueManagerID := "b8e1aeda078009cf3db74e90d5d42328"
				getQueueManagerAvailableUpgradeVersionsOptionsModel := mqcloudService.NewGetQueueManagerAvailableUpgradeVersionsOptions(serviceInstanceGuid, queueManagerID)
				getQueueManagerAvailableUpgradeVersionsOptionsModel.SetServiceInstanceGuid("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				getQueueManagerAvailableUpgradeVersionsOptionsModel.SetQueueManagerID("b8e1aeda078009cf3db74e90d5d42328")
				getQueueManagerAvailableUpgradeVersionsOptionsModel.SetHeaders(map[string]string{"foo": "bar"})
				Expect(getQueueManagerAvailableUpgradeVersionsOptionsModel).ToNot(BeNil())
				Expect(getQueueManagerAvailableUpgradeVersionsOptionsModel.ServiceInstanceGuid).To(Equal(core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")))
				Expect(getQueueManagerAvailableUpgradeVersionsOptionsModel.QueueManagerID).To(Equal(core.StringPtr("b8e1aeda078009cf3db74e90d5d42328")))
				Expect(getQueueManagerAvailableUpgradeVersionsOptionsModel.Headers).To(Equal(map[string]string{"foo": "bar"}))
			})
			It(`Invoke NewGetQueueManagerConnectionInfoOptions successfully`, func() {
				// Construct an instance of the GetQueueManagerConnectionInfoOptions model
				serviceInstanceGuid := "a2b4d4bc-dadb-4637-bcec-9b7d1e723af8"
				queueManagerID := "b8e1aeda078009cf3db74e90d5d42328"
				getQueueManagerConnectionInfoOptionsModel := mqcloudService.NewGetQueueManagerConnectionInfoOptions(serviceInstanceGuid, queueManagerID)
				getQueueManagerConnectionInfoOptionsModel.SetServiceInstanceGuid("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				getQueueManagerConnectionInfoOptionsModel.SetQueueManagerID("b8e1aeda078009cf3db74e90d5d42328")
				getQueueManagerConnectionInfoOptionsModel.SetHeaders(map[string]string{"foo": "bar"})
				Expect(getQueueManagerConnectionInfoOptionsModel).ToNot(BeNil())
				Expect(getQueueManagerConnectionInfoOptionsModel.ServiceInstanceGuid).To(Equal(core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")))
				Expect(getQueueManagerConnectionInfoOptionsModel.QueueManagerID).To(Equal(core.StringPtr("b8e1aeda078009cf3db74e90d5d42328")))
				Expect(getQueueManagerConnectionInfoOptionsModel.Headers).To(Equal(map[string]string{"foo": "bar"}))
			})
			It(`Invoke NewGetQueueManagerOptions successfully`, func() {
				// Construct an instance of the GetQueueManagerOptions model
				serviceInstanceGuid := "a2b4d4bc-dadb-4637-bcec-9b7d1e723af8"
				queueManagerID := "b8e1aeda078009cf3db74e90d5d42328"
				getQueueManagerOptionsModel := mqcloudService.NewGetQueueManagerOptions(serviceInstanceGuid, queueManagerID)
				getQueueManagerOptionsModel.SetServiceInstanceGuid("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				getQueueManagerOptionsModel.SetQueueManagerID("b8e1aeda078009cf3db74e90d5d42328")
				getQueueManagerOptionsModel.SetHeaders(map[string]string{"foo": "bar"})
				Expect(getQueueManagerOptionsModel).ToNot(BeNil())
				Expect(getQueueManagerOptionsModel.ServiceInstanceGuid).To(Equal(core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")))
				Expect(getQueueManagerOptionsModel.QueueManagerID).To(Equal(core.StringPtr("b8e1aeda078009cf3db74e90d5d42328")))
				Expect(getQueueManagerOptionsModel.Headers).To(Equal(map[string]string{"foo": "bar"}))
			})
			It(`Invoke NewGetQueueManagerStatusOptions successfully`, func() {
				// Construct an instance of the GetQueueManagerStatusOptions model
				serviceInstanceGuid := "a2b4d4bc-dadb-4637-bcec-9b7d1e723af8"
				queueManagerID := "b8e1aeda078009cf3db74e90d5d42328"
				getQueueManagerStatusOptionsModel := mqcloudService.NewGetQueueManagerStatusOptions(serviceInstanceGuid, queueManagerID)
				getQueueManagerStatusOptionsModel.SetServiceInstanceGuid("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				getQueueManagerStatusOptionsModel.SetQueueManagerID("b8e1aeda078009cf3db74e90d5d42328")
				getQueueManagerStatusOptionsModel.SetHeaders(map[string]string{"foo": "bar"})
				Expect(getQueueManagerStatusOptionsModel).ToNot(BeNil())
				Expect(getQueueManagerStatusOptionsModel.ServiceInstanceGuid).To(Equal(core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")))
				Expect(getQueueManagerStatusOptionsModel.QueueManagerID).To(Equal(core.StringPtr("b8e1aeda078009cf3db74e90d5d42328")))
				Expect(getQueueManagerStatusOptionsModel.Headers).To(Equal(map[string]string{"foo": "bar"}))
			})
			It(`Invoke NewGetTrustStoreCertificateOptions successfully`, func() {
				// Construct an instance of the GetTrustStoreCertificateOptions model
				serviceInstanceGuid := "a2b4d4bc-dadb-4637-bcec-9b7d1e723af8"
				queueManagerID := "b8e1aeda078009cf3db74e90d5d42328"
				certificateID := "9b7d1e723af8233"
				getTrustStoreCertificateOptionsModel := mqcloudService.NewGetTrustStoreCertificateOptions(serviceInstanceGuid, queueManagerID, certificateID)
				getTrustStoreCertificateOptionsModel.SetServiceInstanceGuid("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				getTrustStoreCertificateOptionsModel.SetQueueManagerID("b8e1aeda078009cf3db74e90d5d42328")
				getTrustStoreCertificateOptionsModel.SetCertificateID("9b7d1e723af8233")
				getTrustStoreCertificateOptionsModel.SetHeaders(map[string]string{"foo": "bar"})
				Expect(getTrustStoreCertificateOptionsModel).ToNot(BeNil())
				Expect(getTrustStoreCertificateOptionsModel.ServiceInstanceGuid).To(Equal(core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")))
				Expect(getTrustStoreCertificateOptionsModel.QueueManagerID).To(Equal(core.StringPtr("b8e1aeda078009cf3db74e90d5d42328")))
				Expect(getTrustStoreCertificateOptionsModel.CertificateID).To(Equal(core.StringPtr("9b7d1e723af8233")))
				Expect(getTrustStoreCertificateOptionsModel.Headers).To(Equal(map[string]string{"foo": "bar"}))
			})
			It(`Invoke NewGetUsageDetailsOptions successfully`, func() {
				// Construct an instance of the GetUsageDetailsOptions model
				serviceInstanceGuid := "a2b4d4bc-dadb-4637-bcec-9b7d1e723af8"
				getUsageDetailsOptionsModel := mqcloudService.NewGetUsageDetailsOptions(serviceInstanceGuid)
				getUsageDetailsOptionsModel.SetServiceInstanceGuid("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				getUsageDetailsOptionsModel.SetHeaders(map[string]string{"foo": "bar"})
				Expect(getUsageDetailsOptionsModel).ToNot(BeNil())
				Expect(getUsageDetailsOptionsModel.ServiceInstanceGuid).To(Equal(core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")))
				Expect(getUsageDetailsOptionsModel.Headers).To(Equal(map[string]string{"foo": "bar"}))
			})
			It(`Invoke NewGetUserOptions successfully`, func() {
				// Construct an instance of the GetUserOptions model
				serviceInstanceGuid := "a2b4d4bc-dadb-4637-bcec-9b7d1e723af8"
				userID := "31a413dd84346effc8895b6ba4641641"
				getUserOptionsModel := mqcloudService.NewGetUserOptions(serviceInstanceGuid, userID)
				getUserOptionsModel.SetServiceInstanceGuid("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				getUserOptionsModel.SetUserID("31a413dd84346effc8895b6ba4641641")
				getUserOptionsModel.SetHeaders(map[string]string{"foo": "bar"})
				Expect(getUserOptionsModel).ToNot(BeNil())
				Expect(getUserOptionsModel.ServiceInstanceGuid).To(Equal(core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")))
				Expect(getUserOptionsModel.UserID).To(Equal(core.StringPtr("31a413dd84346effc8895b6ba4641641")))
				Expect(getUserOptionsModel.Headers).To(Equal(map[string]string{"foo": "bar"}))
			})
			It(`Invoke NewGetVirtualPrivateEndpointGatewayOptions successfully`, func() {
				// Construct an instance of the GetVirtualPrivateEndpointGatewayOptions model
				serviceInstanceGuid := "a2b4d4bc-dadb-4637-bcec-9b7d1e723af8"
				virtualPrivateEndpointGatewayGuid := "r010-ebab3c08-c9a8-40c4-8869-61c09ddf7b44"
				getVirtualPrivateEndpointGatewayOptionsModel := mqcloudService.NewGetVirtualPrivateEndpointGatewayOptions(serviceInstanceGuid, virtualPrivateEndpointGatewayGuid)
				getVirtualPrivateEndpointGatewayOptionsModel.SetServiceInstanceGuid("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				getVirtualPrivateEndpointGatewayOptionsModel.SetVirtualPrivateEndpointGatewayGuid("r010-ebab3c08-c9a8-40c4-8869-61c09ddf7b44")
				getVirtualPrivateEndpointGatewayOptionsModel.SetTrustedProfile("crn:v1:bluemix:public:iam-identity::a/5d5ff2a9001c4055ab1408e9bf185f48::profile:Profile-1c0a8609-ca25-4ad2-a09b-aea472d34afc")
				getVirtualPrivateEndpointGatewayOptionsModel.SetHeaders(map[string]string{"foo": "bar"})
				Expect(getVirtualPrivateEndpointGatewayOptionsModel).ToNot(BeNil())
				Expect(getVirtualPrivateEndpointGatewayOptionsModel.ServiceInstanceGuid).To(Equal(core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")))
				Expect(getVirtualPrivateEndpointGatewayOptionsModel.VirtualPrivateEndpointGatewayGuid).To(Equal(core.StringPtr("r010-ebab3c08-c9a8-40c4-8869-61c09ddf7b44")))
				Expect(getVirtualPrivateEndpointGatewayOptionsModel.TrustedProfile).To(Equal(core.StringPtr("crn:v1:bluemix:public:iam-identity::a/5d5ff2a9001c4055ab1408e9bf185f48::profile:Profile-1c0a8609-ca25-4ad2-a09b-aea472d34afc")))
				Expect(getVirtualPrivateEndpointGatewayOptionsModel.Headers).To(Equal(map[string]string{"foo": "bar"}))
			})
			It(`Invoke NewListApplicationsOptions successfully`, func() {
				// Construct an instance of the ListApplicationsOptions model
				serviceInstanceGuid := "a2b4d4bc-dadb-4637-bcec-9b7d1e723af8"
				listApplicationsOptionsModel := mqcloudService.NewListApplicationsOptions(serviceInstanceGuid)
				listApplicationsOptionsModel.SetServiceInstanceGuid("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				listApplicationsOptionsModel.SetOffset(int64(0))
				listApplicationsOptionsModel.SetLimit(int64(10))
				listApplicationsOptionsModel.SetHeaders(map[string]string{"foo": "bar"})
				Expect(listApplicationsOptionsModel).ToNot(BeNil())
				Expect(listApplicationsOptionsModel.ServiceInstanceGuid).To(Equal(core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")))
				Expect(listApplicationsOptionsModel.Offset).To(Equal(core.Int64Ptr(int64(0))))
				Expect(listApplicationsOptionsModel.Limit).To(Equal(core.Int64Ptr(int64(10))))
				Expect(listApplicationsOptionsModel.Headers).To(Equal(map[string]string{"foo": "bar"}))
			})
			It(`Invoke NewListKeyStoreCertificatesOptions successfully`, func() {
				// Construct an instance of the ListKeyStoreCertificatesOptions model
				serviceInstanceGuid := "a2b4d4bc-dadb-4637-bcec-9b7d1e723af8"
				queueManagerID := "b8e1aeda078009cf3db74e90d5d42328"
				listKeyStoreCertificatesOptionsModel := mqcloudService.NewListKeyStoreCertificatesOptions(serviceInstanceGuid, queueManagerID)
				listKeyStoreCertificatesOptionsModel.SetServiceInstanceGuid("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				listKeyStoreCertificatesOptionsModel.SetQueueManagerID("b8e1aeda078009cf3db74e90d5d42328")
				listKeyStoreCertificatesOptionsModel.SetHeaders(map[string]string{"foo": "bar"})
				Expect(listKeyStoreCertificatesOptionsModel).ToNot(BeNil())
				Expect(listKeyStoreCertificatesOptionsModel.ServiceInstanceGuid).To(Equal(core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")))
				Expect(listKeyStoreCertificatesOptionsModel.QueueManagerID).To(Equal(core.StringPtr("b8e1aeda078009cf3db74e90d5d42328")))
				Expect(listKeyStoreCertificatesOptionsModel.Headers).To(Equal(map[string]string{"foo": "bar"}))
			})
			It(`Invoke NewListQueueManagersOptions successfully`, func() {
				// Construct an instance of the ListQueueManagersOptions model
				serviceInstanceGuid := "a2b4d4bc-dadb-4637-bcec-9b7d1e723af8"
				listQueueManagersOptionsModel := mqcloudService.NewListQueueManagersOptions(serviceInstanceGuid)
				listQueueManagersOptionsModel.SetServiceInstanceGuid("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				listQueueManagersOptionsModel.SetOffset(int64(0))
				listQueueManagersOptionsModel.SetLimit(int64(10))
				listQueueManagersOptionsModel.SetHeaders(map[string]string{"foo": "bar"})
				Expect(listQueueManagersOptionsModel).ToNot(BeNil())
				Expect(listQueueManagersOptionsModel.ServiceInstanceGuid).To(Equal(core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")))
				Expect(listQueueManagersOptionsModel.Offset).To(Equal(core.Int64Ptr(int64(0))))
				Expect(listQueueManagersOptionsModel.Limit).To(Equal(core.Int64Ptr(int64(10))))
				Expect(listQueueManagersOptionsModel.Headers).To(Equal(map[string]string{"foo": "bar"}))
			})
			It(`Invoke NewListTrustStoreCertificatesOptions successfully`, func() {
				// Construct an instance of the ListTrustStoreCertificatesOptions model
				serviceInstanceGuid := "a2b4d4bc-dadb-4637-bcec-9b7d1e723af8"
				queueManagerID := "b8e1aeda078009cf3db74e90d5d42328"
				listTrustStoreCertificatesOptionsModel := mqcloudService.NewListTrustStoreCertificatesOptions(serviceInstanceGuid, queueManagerID)
				listTrustStoreCertificatesOptionsModel.SetServiceInstanceGuid("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				listTrustStoreCertificatesOptionsModel.SetQueueManagerID("b8e1aeda078009cf3db74e90d5d42328")
				listTrustStoreCertificatesOptionsModel.SetHeaders(map[string]string{"foo": "bar"})
				Expect(listTrustStoreCertificatesOptionsModel).ToNot(BeNil())
				Expect(listTrustStoreCertificatesOptionsModel.ServiceInstanceGuid).To(Equal(core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")))
				Expect(listTrustStoreCertificatesOptionsModel.QueueManagerID).To(Equal(core.StringPtr("b8e1aeda078009cf3db74e90d5d42328")))
				Expect(listTrustStoreCertificatesOptionsModel.Headers).To(Equal(map[string]string{"foo": "bar"}))
			})
			It(`Invoke NewListUsersOptions successfully`, func() {
				// Construct an instance of the ListUsersOptions model
				serviceInstanceGuid := "a2b4d4bc-dadb-4637-bcec-9b7d1e723af8"
				listUsersOptionsModel := mqcloudService.NewListUsersOptions(serviceInstanceGuid)
				listUsersOptionsModel.SetServiceInstanceGuid("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				listUsersOptionsModel.SetOffset(int64(0))
				listUsersOptionsModel.SetLimit(int64(10))
				listUsersOptionsModel.SetHeaders(map[string]string{"foo": "bar"})
				Expect(listUsersOptionsModel).ToNot(BeNil())
				Expect(listUsersOptionsModel.ServiceInstanceGuid).To(Equal(core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")))
				Expect(listUsersOptionsModel.Offset).To(Equal(core.Int64Ptr(int64(0))))
				Expect(listUsersOptionsModel.Limit).To(Equal(core.Int64Ptr(int64(10))))
				Expect(listUsersOptionsModel.Headers).To(Equal(map[string]string{"foo": "bar"}))
			})
			It(`Invoke NewListVirtualPrivateEndpointGatewaysOptions successfully`, func() {
				// Construct an instance of the ListVirtualPrivateEndpointGatewaysOptions model
				serviceInstanceGuid := "a2b4d4bc-dadb-4637-bcec-9b7d1e723af8"
				listVirtualPrivateEndpointGatewaysOptionsModel := mqcloudService.NewListVirtualPrivateEndpointGatewaysOptions(serviceInstanceGuid)
				listVirtualPrivateEndpointGatewaysOptionsModel.SetServiceInstanceGuid("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				listVirtualPrivateEndpointGatewaysOptionsModel.SetTrustedProfile("crn:v1:bluemix:public:iam-identity::a/5d5ff2a9001c4055ab1408e9bf185f48::profile:Profile-1c0a8609-ca25-4ad2-a09b-aea472d34afc")
				listVirtualPrivateEndpointGatewaysOptionsModel.SetStart("r010-ebab3c08-c9a8-40c4-8869-61c09ddf7b44")
				listVirtualPrivateEndpointGatewaysOptionsModel.SetLimit(int64(10))
				listVirtualPrivateEndpointGatewaysOptionsModel.SetHeaders(map[string]string{"foo": "bar"})
				Expect(listVirtualPrivateEndpointGatewaysOptionsModel).ToNot(BeNil())
				Expect(listVirtualPrivateEndpointGatewaysOptionsModel.ServiceInstanceGuid).To(Equal(core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")))
				Expect(listVirtualPrivateEndpointGatewaysOptionsModel.TrustedProfile).To(Equal(core.StringPtr("crn:v1:bluemix:public:iam-identity::a/5d5ff2a9001c4055ab1408e9bf185f48::profile:Profile-1c0a8609-ca25-4ad2-a09b-aea472d34afc")))
				Expect(listVirtualPrivateEndpointGatewaysOptionsModel.Start).To(Equal(core.StringPtr("r010-ebab3c08-c9a8-40c4-8869-61c09ddf7b44")))
				Expect(listVirtualPrivateEndpointGatewaysOptionsModel.Limit).To(Equal(core.Int64Ptr(int64(10))))
				Expect(listVirtualPrivateEndpointGatewaysOptionsModel.Headers).To(Equal(map[string]string{"foo": "bar"}))
			})
			It(`Invoke NewSetApplicationNameOptions successfully`, func() {
				// Construct an instance of the SetApplicationNameOptions model
				serviceInstanceGuid := "a2b4d4bc-dadb-4637-bcec-9b7d1e723af8"
				applicationID := "0123456789ABCDEF0123456789ABCDEF"
				setApplicationNameOptionsName := "testString"
				setApplicationNameOptionsModel := mqcloudService.NewSetApplicationNameOptions(serviceInstanceGuid, applicationID, setApplicationNameOptionsName)
				setApplicationNameOptionsModel.SetServiceInstanceGuid("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				setApplicationNameOptionsModel.SetApplicationID("0123456789ABCDEF0123456789ABCDEF")
				setApplicationNameOptionsModel.SetName("testString")
				setApplicationNameOptionsModel.SetHeaders(map[string]string{"foo": "bar"})
				Expect(setApplicationNameOptionsModel).ToNot(BeNil())
				Expect(setApplicationNameOptionsModel.ServiceInstanceGuid).To(Equal(core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")))
				Expect(setApplicationNameOptionsModel.ApplicationID).To(Equal(core.StringPtr("0123456789ABCDEF0123456789ABCDEF")))
				Expect(setApplicationNameOptionsModel.Name).To(Equal(core.StringPtr("testString")))
				Expect(setApplicationNameOptionsModel.Headers).To(Equal(map[string]string{"foo": "bar"}))
			})
			It(`Invoke NewSetCertificateAmsChannelsOptions successfully`, func() {
				// Construct an instance of the ChannelDetails model
				channelDetailsModel := new(mqcloudv1.ChannelDetails)
				Expect(channelDetailsModel).ToNot(BeNil())
				channelDetailsModel.Name = core.StringPtr("testString")
				Expect(channelDetailsModel.Name).To(Equal(core.StringPtr("testString")))

				// Construct an instance of the SetCertificateAmsChannelsOptions model
				queueManagerID := "b8e1aeda078009cf3db74e90d5d42328"
				certificateID := "9b7d1e723af8233"
				serviceInstanceGuid := "a2b4d4bc-dadb-4637-bcec-9b7d1e723af8"
				setCertificateAmsChannelsOptionsChannels := []mqcloudv1.ChannelDetails{}
				setCertificateAmsChannelsOptionsModel := mqcloudService.NewSetCertificateAmsChannelsOptions(queueManagerID, certificateID, serviceInstanceGuid, setCertificateAmsChannelsOptionsChannels)
				setCertificateAmsChannelsOptionsModel.SetQueueManagerID("b8e1aeda078009cf3db74e90d5d42328")
				setCertificateAmsChannelsOptionsModel.SetCertificateID("9b7d1e723af8233")
				setCertificateAmsChannelsOptionsModel.SetServiceInstanceGuid("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				setCertificateAmsChannelsOptionsModel.SetChannels([]mqcloudv1.ChannelDetails{*channelDetailsModel})
				setCertificateAmsChannelsOptionsModel.SetUpdateStrategy("replace")
				setCertificateAmsChannelsOptionsModel.SetHeaders(map[string]string{"foo": "bar"})
				Expect(setCertificateAmsChannelsOptionsModel).ToNot(BeNil())
				Expect(setCertificateAmsChannelsOptionsModel.QueueManagerID).To(Equal(core.StringPtr("b8e1aeda078009cf3db74e90d5d42328")))
				Expect(setCertificateAmsChannelsOptionsModel.CertificateID).To(Equal(core.StringPtr("9b7d1e723af8233")))
				Expect(setCertificateAmsChannelsOptionsModel.ServiceInstanceGuid).To(Equal(core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")))
				Expect(setCertificateAmsChannelsOptionsModel.Channels).To(Equal([]mqcloudv1.ChannelDetails{*channelDetailsModel}))
				Expect(setCertificateAmsChannelsOptionsModel.UpdateStrategy).To(Equal(core.StringPtr("replace")))
				Expect(setCertificateAmsChannelsOptionsModel.Headers).To(Equal(map[string]string{"foo": "bar"}))
			})
			It(`Invoke NewSetQueueManagerVersionOptions successfully`, func() {
				// Construct an instance of the SetQueueManagerVersionOptions model
				serviceInstanceGuid := "a2b4d4bc-dadb-4637-bcec-9b7d1e723af8"
				queueManagerID := "b8e1aeda078009cf3db74e90d5d42328"
				setQueueManagerVersionOptionsVersion := "9.3.2_2"
				setQueueManagerVersionOptionsModel := mqcloudService.NewSetQueueManagerVersionOptions(serviceInstanceGuid, queueManagerID, setQueueManagerVersionOptionsVersion)
				setQueueManagerVersionOptionsModel.SetServiceInstanceGuid("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				setQueueManagerVersionOptionsModel.SetQueueManagerID("b8e1aeda078009cf3db74e90d5d42328")
				setQueueManagerVersionOptionsModel.SetVersion("9.3.2_2")
				setQueueManagerVersionOptionsModel.SetHeaders(map[string]string{"foo": "bar"})
				Expect(setQueueManagerVersionOptionsModel).ToNot(BeNil())
				Expect(setQueueManagerVersionOptionsModel.ServiceInstanceGuid).To(Equal(core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")))
				Expect(setQueueManagerVersionOptionsModel.QueueManagerID).To(Equal(core.StringPtr("b8e1aeda078009cf3db74e90d5d42328")))
				Expect(setQueueManagerVersionOptionsModel.Version).To(Equal(core.StringPtr("9.3.2_2")))
				Expect(setQueueManagerVersionOptionsModel.Headers).To(Equal(map[string]string{"foo": "bar"}))
			})
			It(`Invoke NewSetUserNameOptions successfully`, func() {
				// Construct an instance of the SetUserNameOptions model
				serviceInstanceGuid := "a2b4d4bc-dadb-4637-bcec-9b7d1e723af8"
				userID := "31a413dd84346effc8895b6ba4641641"
				setUserNameOptionsName := "testString"
				setUserNameOptionsModel := mqcloudService.NewSetUserNameOptions(serviceInstanceGuid, userID, setUserNameOptionsName)
				setUserNameOptionsModel.SetServiceInstanceGuid("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")
				setUserNameOptionsModel.SetUserID("31a413dd84346effc8895b6ba4641641")
				setUserNameOptionsModel.SetName("testString")
				setUserNameOptionsModel.SetHeaders(map[string]string{"foo": "bar"})
				Expect(setUserNameOptionsModel).ToNot(BeNil())
				Expect(setUserNameOptionsModel.ServiceInstanceGuid).To(Equal(core.StringPtr("a2b4d4bc-dadb-4637-bcec-9b7d1e723af8")))
				Expect(setUserNameOptionsModel.UserID).To(Equal(core.StringPtr("31a413dd84346effc8895b6ba4641641")))
				Expect(setUserNameOptionsModel.Name).To(Equal(core.StringPtr("testString")))
				Expect(setUserNameOptionsModel.Headers).To(Equal(map[string]string{"foo": "bar"}))
			})
		})
	})
	Describe(`Model unmarshaling tests`, func() {
		It(`Invoke UnmarshalChannelDetails successfully`, func() {
			// Construct an instance of the model.
			model := new(mqcloudv1.ChannelDetails)
			model.Name = core.StringPtr("testString")

			b, err := json.Marshal(model)
			Expect(err).To(BeNil())

			var raw map[string]json.RawMessage
			err = json.Unmarshal(b, &raw)
			Expect(err).To(BeNil())

			var result *mqcloudv1.ChannelDetails
			err = mqcloudv1.UnmarshalChannelDetails(raw, &result)
			Expect(err).To(BeNil())
			Expect(result).ToNot(BeNil())
			Expect(result).To(Equal(model))
		})
	})
	Describe(`Utility function tests`, func() {
		It(`Invoke CreateMockByteArray() successfully`, func() {
			mockByteArray := CreateMockByteArray("VGhpcyBpcyBhIHRlc3Qgb2YgdGhlIGVtZXJnZW5jeSBicm9hZGNhc3Qgc3lzdGVt")
			Expect(mockByteArray).ToNot(BeNil())
		})
		It(`Invoke CreateMockUUID() successfully`, func() {
			mockUUID := CreateMockUUID("9fab83da-98cb-4f18-a7ba-b6f0435c9673")
			Expect(mockUUID).ToNot(BeNil())
		})
		It(`Invoke CreateMockReader() successfully`, func() {
			mockReader := CreateMockReader("This is a test.")
			Expect(mockReader).ToNot(BeNil())
		})
		It(`Invoke CreateMockDate() successfully`, func() {
			mockDate := CreateMockDate("2019-01-01")
			Expect(mockDate).ToNot(BeNil())
		})
		It(`Invoke CreateMockDateTime() successfully`, func() {
			mockDateTime := CreateMockDateTime("2019-01-01T12:00:00.000Z")
			Expect(mockDateTime).ToNot(BeNil())
		})
	})
})

//
// Utility functions used by the generated test code
//

func CreateMockByteArray(encodedString string) *[]byte {
	ba, err := base64.StdEncoding.DecodeString(encodedString)
	if err != nil {
		panic(err)
	}
	return &ba
}

func CreateMockUUID(mockData string) *strfmt.UUID {
	uuid := strfmt.UUID(mockData)
	return &uuid
}

func CreateMockReader(mockData string) io.ReadCloser {
	return io.NopCloser(bytes.NewReader([]byte(mockData)))
}

func CreateMockDate(mockData string) *strfmt.Date {
	d, err := core.ParseDate(mockData)
	if err != nil {
		return nil
	}
	return &d
}

func CreateMockDateTime(mockData string) *strfmt.DateTime {
	d, err := core.ParseDateTime(mockData)
	if err != nil {
		return nil
	}
	return &d
}

func SetTestEnvironment(testEnvironment map[string]string) {
	for key, value := range testEnvironment {
		os.Setenv(key, value)
	}
}

func ClearTestEnvironment(testEnvironment map[string]string) {
	for key := range testEnvironment {
		os.Unsetenv(key)
	}
}
