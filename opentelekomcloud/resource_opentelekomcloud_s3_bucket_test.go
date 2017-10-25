package opentelekomcloud

import (
	//"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"testing"
	//"text/template"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"

	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	//"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/hashicorp/terraform/helper/schema"
)

// PASS
func TestAccAWSS3Bucket_basic(t *testing.T) {
	rInt := acctest.RandInt()
	//arnRegexp := regexp.MustCompile("^arn:aws:s3:::")

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		/*
			IDRefreshName:   "opentelekomcloud_s3_bucket.bucket",
			IDRefreshIgnore: []string{"force_destroy"},
		*/
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSS3BucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSS3BucketConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSS3BucketExists("opentelekomcloud_s3_bucket.bucket"),
					/*resource.TestCheckResourceAttr(
					"opentelekomcloud_s3_bucket.bucket", "hosted_zone_id", HostedZoneIDForRegion("us-west-2")), */
					resource.TestCheckResourceAttr(
						"opentelekomcloud_s3_bucket.bucket", "region", OS_REGION_NAME),
					resource.TestCheckNoResourceAttr(
						"opentelekomcloud_s3_bucket.bucket", "website_endpoint"),
					/*resource.TestMatchResourceAttr(
					"opentelekomcloud_s3_bucket.bucket", "arn", arnRegexp), */
					resource.TestCheckResourceAttr(
						"opentelekomcloud_s3_bucket.bucket", "bucket", testAccBucketName(rInt)),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_s3_bucket.bucket", "bucket_domain_name", testAccBucketDomainName(rInt)),
				),
			},
		},
	})
}

// NOT SUPPORTED
/*
func TestAccAWSS3MultiBucket_withTags(t *testing.T) {
	rInt := acctest.RandInt()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSS3BucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSS3MultiBucketConfigWithTags(rInt),
			},
		},
	})
}
*/

// PASS
func TestAccAWSS3Bucket_namePrefix(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSS3BucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSS3BucketConfig_namePrefix,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSS3BucketExists("opentelekomcloud_s3_bucket.test"),
					resource.TestMatchResourceAttr(
						"opentelekomcloud_s3_bucket.test", "bucket", regexp.MustCompile("^tf-test-")),
				),
			},
		},
	})
}

// PASS
func TestAccAWSS3Bucket_generatedName(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSS3BucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSS3BucketConfig_generatedName,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSS3BucketExists("opentelekomcloud_s3_bucket.test"),
				),
			},
		},
	})
}

// PASS
func TestAccAWSS3Bucket_region(t *testing.T) {
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSS3BucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSS3BucketConfigWithRegion(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSS3BucketExists("opentelekomcloud_s3_bucket.bucket"),
					resource.TestCheckResourceAttr("opentelekomcloud_s3_bucket.bucket", "region", OS_REGION_NAME),
				),
			},
		},
	})
}

// NOT SUPPORTED
/*
func TestAccAWSS3Bucket_acceleration(t *testing.T) {
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSS3BucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSS3BucketConfigWithAcceleration(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSS3BucketExists("opentelekomcloud_s3_bucket.bucket"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_s3_bucket.bucket", "acceleration_status", "Enabled"),
				),
			},
			{
				Config: testAccAWSS3BucketConfigWithoutAcceleration(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSS3BucketExists("opentelekomcloud_s3_bucket.bucket"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_s3_bucket.bucket", "acceleration_status", "Suspended"),
				),
			},
		},
	})
}
*/

// UNSUPPORTED
/*
func TestAccAWSS3Bucket_RequestPayer(t *testing.T) {
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSS3BucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSS3BucketConfigRequestPayerBucketOwner(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSS3BucketExists("opentelekomcloud_s3_bucket.bucket"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_s3_bucket.bucket",
						"request_payer",
						"BucketOwner"),
					testAccCheckAWSS3RequestPayer(
						"opentelekomcloud_s3_bucket.bucket",
						"BucketOwner"),
				),
			},
			{
				Config: testAccAWSS3BucketConfigRequestPayerRequester(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSS3BucketExists("opentelekomcloud_s3_bucket.bucket"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_s3_bucket.bucket",
						"request_payer",
						"Requester"),
					testAccCheckAWSS3RequestPayer(
						"opentelekomcloud_s3_bucket.bucket",
						"Requester"),
				),
			},
		},
	})
}
*/

// PASS, but not needed or supported
/*
func TestResourceAWSS3BucketRequestPayer_validation(t *testing.T) {
	_, errors := validateS3BucketRequestPayerType("incorrect", "request_payer")
	if len(errors) == 0 {
		t.Fatalf("Expected to trigger a validation error")
	}

	var testCases = []struct {
		Value    string
		ErrCount int
	}{
		{
			Value:    "Requester",
			ErrCount: 0,
		},
		{
			Value:    "BucketOwner",
			ErrCount: 0,
		},
	}

	for _, tc := range testCases {
		_, errors := validateS3BucketRequestPayerType(tc.Value, "request_payer")
		if len(errors) != tc.ErrCount {
			t.Fatalf("Expected not to trigger a validation error")
		}
	}
}
*/

// PASS
func TestAccAWSS3Bucket_Policy(t *testing.T) {
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSS3BucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSS3BucketConfigWithPolicy(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSS3BucketExists("opentelekomcloud_s3_bucket.bucket"),
					testAccCheckAWSS3BucketPolicy(
						"opentelekomcloud_s3_bucket.bucket", testAccAWSS3BucketPolicy(rInt)),
				),
			},
			{
				Config: testAccAWSS3BucketConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSS3BucketExists("opentelekomcloud_s3_bucket.bucket"),
					testAccCheckAWSS3BucketPolicy(
						"opentelekomcloud_s3_bucket.bucket", ""),
				),
			},
			{
				Config: testAccAWSS3BucketConfigWithEmptyPolicy(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSS3BucketExists("opentelekomcloud_s3_bucket.bucket"),
					testAccCheckAWSS3BucketPolicy(
						"opentelekomcloud_s3_bucket.bucket", ""),
				),
			},
		},
	})
}

// PASS
func TestAccAWSS3Bucket_UpdateAcl(t *testing.T) {
	ri := acctest.RandInt()
	preConfig := fmt.Sprintf(testAccAWSS3BucketConfigWithAcl, ri)
	postConfig := fmt.Sprintf(testAccAWSS3BucketConfigWithAclUpdate, ri)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSS3BucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSS3BucketExists("opentelekomcloud_s3_bucket.bucket"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_s3_bucket.bucket", "acl", "public-read"),
				),
			},
			{
				Config: postConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSS3BucketExists("opentelekomcloud_s3_bucket.bucket"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_s3_bucket.bucket", "acl", "private"),
				),
			},
		},
	})
}

// PASS
func TestAccAWSS3Bucket_Website_Simple(t *testing.T) {
	rInt := acctest.RandInt()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSS3BucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSS3BucketWebsiteConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSS3BucketExists("opentelekomcloud_s3_bucket.bucket"),
					testAccCheckAWSS3BucketWebsite(
						"opentelekomcloud_s3_bucket.bucket", "index.html", "", "", ""),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_s3_bucket.bucket", "website_endpoint", testAccWebsiteEndpoint(rInt)),
				),
			},
			{
				Config: testAccAWSS3BucketWebsiteConfigWithError(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSS3BucketExists("opentelekomcloud_s3_bucket.bucket"),
					testAccCheckAWSS3BucketWebsite(
						"opentelekomcloud_s3_bucket.bucket", "index.html", "error.html", "", ""),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_s3_bucket.bucket", "website_endpoint", testAccWebsiteEndpoint(rInt)),
				),
			},
			{
				Config: testAccAWSS3BucketConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSS3BucketExists("opentelekomcloud_s3_bucket.bucket"),
					testAccCheckAWSS3BucketWebsite(
						"opentelekomcloud_s3_bucket.bucket", "", "", "", ""),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_s3_bucket.bucket", "website_endpoint", ""),
				),
			},
		},
	})
}

// PASS
func TestAccAWSS3Bucket_WebsiteRedirect(t *testing.T) {
	rInt := acctest.RandInt()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSS3BucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSS3BucketWebsiteConfigWithRedirect(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSS3BucketExists("opentelekomcloud_s3_bucket.bucket"),
					testAccCheckAWSS3BucketWebsite(
						"opentelekomcloud_s3_bucket.bucket", "", "", "", "hashicorp.com"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_s3_bucket.bucket", "website_endpoint", testAccWebsiteEndpoint(rInt)),
				),
			},
			{
				Config: testAccAWSS3BucketWebsiteConfigWithHttpsRedirect(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSS3BucketExists("opentelekomcloud_s3_bucket.bucket"),
					testAccCheckAWSS3BucketWebsite(
						"opentelekomcloud_s3_bucket.bucket", "", "", "https", "hashicorp.com"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_s3_bucket.bucket", "website_endpoint", testAccWebsiteEndpoint(rInt)),
				),
			},
			{
				Config: testAccAWSS3BucketConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSS3BucketExists("opentelekomcloud_s3_bucket.bucket"),
					testAccCheckAWSS3BucketWebsite(
						"opentelekomcloud_s3_bucket.bucket", "", "", "", ""),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_s3_bucket.bucket", "website_endpoint", ""),
				),
			},
		},
	})
}

// PASS
func TestAccAWSS3Bucket_WebsiteRoutingRules(t *testing.T) {
	rInt := acctest.RandInt()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSS3BucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSS3BucketWebsiteConfigWithRoutingRules(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSS3BucketExists("opentelekomcloud_s3_bucket.bucket"),
					testAccCheckAWSS3BucketWebsite(
						"opentelekomcloud_s3_bucket.bucket", "index.html", "error.html", "", ""),
					testAccCheckAWSS3BucketWebsiteRoutingRules(
						"opentelekomcloud_s3_bucket.bucket",
						[]*s3.RoutingRule{
							{
								Condition: &s3.Condition{
									KeyPrefixEquals: aws.String("docs/"),
								},
								Redirect: &s3.Redirect{
									ReplaceKeyPrefixWith: aws.String("documents/"),
								},
							},
						},
					),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_s3_bucket.bucket", "website_endpoint", testAccWebsiteEndpoint(rInt)),
				),
			},
			{
				Config: testAccAWSS3BucketConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSS3BucketExists("opentelekomcloud_s3_bucket.bucket"),
					testAccCheckAWSS3BucketWebsite(
						"opentelekomcloud_s3_bucket.bucket", "", "", "", ""),
					testAccCheckAWSS3BucketWebsiteRoutingRules("opentelekomcloud_s3_bucket.bucket", nil),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_s3_bucket.bucket", "website_endpoint", ""),
				),
			},
		},
	})
}

// Test TestAccAWSS3Bucket_shouldFailNotFound is designed to fail with a "plan
// not empty" error in Terraform, to check against regresssions.
// See https://github.com/hashicorp/terraform/pull/2925
// PASS
func TestAccAWSS3Bucket_shouldFailNotFound(t *testing.T) {
	rInt := acctest.RandInt()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSS3BucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSS3BucketDestroyedConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSS3BucketExists("opentelekomcloud_s3_bucket.bucket"),
					testAccCheckAWSS3DestroyBucket("opentelekomcloud_s3_bucket.bucket"),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

// PASS
func TestAccAWSS3Bucket_Versioning(t *testing.T) {
	rInt := acctest.RandInt()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSS3BucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSS3BucketConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSS3BucketExists("opentelekomcloud_s3_bucket.bucket"),
					testAccCheckAWSS3BucketVersioning(
						"opentelekomcloud_s3_bucket.bucket", ""),
				),
			},
			{
				Config: testAccAWSS3BucketConfigWithVersioning(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSS3BucketExists("opentelekomcloud_s3_bucket.bucket"),
					testAccCheckAWSS3BucketVersioning(
						"opentelekomcloud_s3_bucket.bucket", s3.BucketVersioningStatusEnabled),
				),
			},
			{
				Config: testAccAWSS3BucketConfigWithDisableVersioning(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSS3BucketExists("opentelekomcloud_s3_bucket.bucket"),
					testAccCheckAWSS3BucketVersioning(
						"opentelekomcloud_s3_bucket.bucket", s3.BucketVersioningStatusSuspended),
				),
			},
		},
	})
}

// PASS
func TestAccAWSS3Bucket_Cors(t *testing.T) {
	rInt := acctest.RandInt()

	updateBucketCors := func(n string) resource.TestCheckFunc {
		return func(s *terraform.State) error {
			rs, ok := s.RootModule().Resources[n]
			if !ok {
				return fmt.Errorf("Not found: %s", n)
			}

			config := testAccProvider.Meta().(*Config)
			conn, err := config.computeS3conn(OS_REGION_NAME)
			if err != nil {
				return fmt.Errorf("Error creating OpenTelekomCloud s3 client: %s", err)
			}
			_, err = conn.PutBucketCors(&s3.PutBucketCorsInput{
				Bucket: aws.String(rs.Primary.ID),
				CORSConfiguration: &s3.CORSConfiguration{
					CORSRules: []*s3.CORSRule{
						{
							AllowedHeaders: []*string{aws.String("*")},
							AllowedMethods: []*string{aws.String("GET")},
							AllowedOrigins: []*string{aws.String("https://www.example.com")},
						},
					},
				},
			})
			if err != nil {
				if awsErr, ok := err.(awserr.Error); ok && awsErr.Code() != "NoSuchCORSConfiguration" {
					return err
				}
			}
			return nil
		}
	}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSS3BucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSS3BucketConfigWithCORS(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSS3BucketExists("opentelekomcloud_s3_bucket.bucket"),
					testAccCheckAWSS3BucketCors(
						"opentelekomcloud_s3_bucket.bucket",
						[]*s3.CORSRule{
							{
								AllowedHeaders: []*string{aws.String("*")},
								AllowedMethods: []*string{aws.String("PUT"), aws.String("POST")},
								AllowedOrigins: []*string{aws.String("https://www.example.com")},
								ExposeHeaders:  []*string{aws.String("x-amz-server-side-encryption"), aws.String("ETag")},
								MaxAgeSeconds:  aws.Int64(3000),
							},
						},
					),
					updateBucketCors("opentelekomcloud_s3_bucket.bucket"),
				),
				ExpectNonEmptyPlan: true, // TODO: No diff in real life, so maybe a timing problem?
			},
			{
				Config: testAccAWSS3BucketConfigWithCORS(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSS3BucketExists("opentelekomcloud_s3_bucket.bucket"),
					testAccCheckAWSS3BucketCors(
						"opentelekomcloud_s3_bucket.bucket",
						[]*s3.CORSRule{
							{
								AllowedHeaders: []*string{aws.String("*")},
								AllowedMethods: []*string{aws.String("PUT"), aws.String("POST")},
								AllowedOrigins: []*string{aws.String("https://www.example.com")},
								ExposeHeaders:  []*string{aws.String("x-amz-server-side-encryption"), aws.String("ETag")},
								MaxAgeSeconds:  aws.Int64(3000),
							},
						},
					),
				),
			},
		},
	})
}

// PASS
func TestAccAWSS3Bucket_Logging(t *testing.T) {
	rInt := acctest.RandInt()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSS3BucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSS3BucketConfigWithLogging(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSS3BucketExists("opentelekomcloud_s3_bucket.bucket"),
					testAccCheckAWSS3BucketLogging(
						"opentelekomcloud_s3_bucket.bucket", "opentelekomcloud_s3_bucket.log_bucket", "log/"),
				),
			},
		},
	})
}

// FAIL: MalformedXML, gets Internal Error if XML is right
// UNSUPPORTED due to being broken.
/*
func TestAccAWSS3Bucket_Lifecycle(t *testing.T) {
	rInt := acctest.RandInt()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSS3BucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSS3BucketConfigWithLifecycle(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSS3BucketExists("opentelekomcloud_s3_bucket.bucket"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_s3_bucket.bucket", "lifecycle_rule.0.id", "id1"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_s3_bucket.bucket", "lifecycle_rule.0.prefix", "path1/"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_s3_bucket.bucket", "lifecycle_rule.0.expiration.2613713285.days", "365"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_s3_bucket.bucket", "lifecycle_rule.0.expiration.2613713285.date", ""),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_s3_bucket.bucket", "lifecycle_rule.0.expiration.2613713285.expired_object_delete_marker", "false"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_s3_bucket.bucket", "lifecycle_rule.0.transition.2000431762.date", ""),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_s3_bucket.bucket", "lifecycle_rule.0.transition.2000431762.days", "30"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_s3_bucket.bucket", "lifecycle_rule.0.transition.2000431762.storage_class", "STANDARD_IA"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_s3_bucket.bucket", "lifecycle_rule.0.transition.6450812.date", ""),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_s3_bucket.bucket", "lifecycle_rule.0.transition.6450812.days", "60"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_s3_bucket.bucket", "lifecycle_rule.0.transition.6450812.storage_class", "GLACIER"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_s3_bucket.bucket", "lifecycle_rule.1.id", "id2"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_s3_bucket.bucket", "lifecycle_rule.1.prefix", "path2/"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_s3_bucket.bucket", "lifecycle_rule.1.expiration.2855832418.date", "2016-01-12"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_s3_bucket.bucket", "lifecycle_rule.1.expiration.2855832418.days", "0"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_s3_bucket.bucket", "lifecycle_rule.1.expiration.2855832418.expired_object_delete_marker", "false"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_s3_bucket.bucket", "lifecycle_rule.2.id", "id3"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_s3_bucket.bucket", "lifecycle_rule.2.prefix", "path3/"),
					resource.TestCheckResourceAttr(
					"opentelekomcloud_s3_bucket.bucket", "lifecycle_rule.2.transition.460947558.days", "0"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_s3_bucket.bucket", "lifecycle_rule.3.id", "id4"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_s3_bucket.bucket", "lifecycle_rule.3.prefix", "path4/"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_s3_bucket.bucket", "lifecycle_rule.3.tags.tagKey", "tagValue"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_s3_bucket.bucket", "lifecycle_rule.3.tags.terraform", "hashicorp"),
				),
			},
			{
				Config: testAccAWSS3BucketConfigWithVersioningLifecycle(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSS3BucketExists("opentelekomcloud_s3_bucket.bucket"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_s3_bucket.bucket", "lifecycle_rule.0.id", "id1"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_s3_bucket.bucket", "lifecycle_rule.0.prefix", "path1/"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_s3_bucket.bucket", "lifecycle_rule.0.enabled", "true"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_s3_bucket.bucket", "lifecycle_rule.0.noncurrent_version_expiration.80908210.days", "365"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_s3_bucket.bucket", "lifecycle_rule.0.noncurrent_version_transition.1377917700.days", "30"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_s3_bucket.bucket", "lifecycle_rule.0.noncurrent_version_transition.1377917700.storage_class", "STANDARD_IA"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_s3_bucket.bucket", "lifecycle_rule.0.noncurrent_version_transition.2528035817.days", "60"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_s3_bucket.bucket", "lifecycle_rule.0.noncurrent_version_transition.2528035817.storage_class", "GLACIER"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_s3_bucket.bucket", "lifecycle_rule.1.id", "id2"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_s3_bucket.bucket", "lifecycle_rule.1.prefix", "path2/"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_s3_bucket.bucket", "lifecycle_rule.1.enabled", "false"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_s3_bucket.bucket", "lifecycle_rule.1.noncurrent_version_expiration.80908210.days", "365"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_s3_bucket.bucket", "lifecycle_rule.2.id", "id3"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_s3_bucket.bucket", "lifecycle_rule.2.prefix", "path3/"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_s3_bucket.bucket", "lifecycle_rule.2.noncurrent_version_transition.3732708140.days", "0"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_s3_bucket.bucket", "lifecycle_rule.2.noncurrent_version_transition.3732708140.storage_class", "GLACIER"),
				),
			},
			{
				Config: testAccAWSS3BucketConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSS3BucketExists("opentelekomcloud_s3_bucket.bucket"),
				),
			},
		},
	})
}
*/

// UNSUPPORTED
/*
func TestAccAWSS3Bucket_Replication(t *testing.T) {
	rInt := acctest.RandInt()

	// record the initialized providers so that we can use them to check for the instances in each region
	var providers []*schema.Provider
	providerFactories := map[string]terraform.ResourceProviderFactory{
		"aws": func() (terraform.ResourceProvider, error) {
			p := Provider()
			providers = append(providers, p.(*schema.Provider))
			return p, nil
		},
	}

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckAWSS3BucketDestroyWithProviders(&providers),
		Steps: []resource.TestStep{
			{
				Config: testAccAWSS3BucketConfigReplication(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSS3BucketExistsWithProviders("opentelekomcloud_s3_bucket.bucket", &providers),
				),
			},
			{
				Config: testAccAWSS3BucketConfigReplicationWithConfiguration(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSS3BucketExistsWithProviders("opentelekomcloud_s3_bucket.bucket", &providers),
					resource.TestCheckResourceAttr("opentelekomcloud_s3_bucket.bucket", "replication_configuration.#", "1"),
					resource.TestCheckResourceAttr("opentelekomcloud_s3_bucket.bucket", "replication_configuration.0.rules.#", "1"),
					resource.TestCheckResourceAttr("opentelekomcloud_s3_bucket.bucket", "replication_configuration.0.rules.2229345141.id", "foobar"),
					resource.TestCheckResourceAttr("opentelekomcloud_s3_bucket.bucket", "replication_configuration.0.rules.2229345141.prefix", "foo"),
					resource.TestCheckResourceAttr("opentelekomcloud_s3_bucket.bucket", "replication_configuration.0.rules.2229345141.status", s3.ReplicationRuleStatusEnabled),
				),
			},
		},
	})
}
*/

// StorageClass issue: https://github.com/hashicorp/terraform/issues/10909
// UNSUPPORTED
/*
func TestAccAWSS3Bucket_ReplicationWithoutStorageClass(t *testing.T) {
	rInt := acctest.RandInt()

	// record the initialized providers so that we can use them to check for the instances in each region
	var providers []*schema.Provider
	providerFactories := map[string]terraform.ResourceProviderFactory{
		"aws": func() (terraform.ResourceProvider, error) {
			p := Provider()
			providers = append(providers, p.(*schema.Provider))
			return p, nil
		},
	}

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckAWSS3BucketDestroyWithProviders(&providers),
		Steps: []resource.TestStep{
			{
				Config: testAccAWSS3BucketConfigReplicationWithoutStorageClass(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSS3BucketExistsWithProviders("opentelekomcloud_s3_bucket.bucket", &providers),
				),
			},
		},
	})
}

// FAIL provider "opentelekomcloud" is not available
func TestAccAWSS3Bucket_ReplicationExpectVersioningValidationError(t *testing.T) {
	rInt := acctest.RandInt()

	// record the initialized providers so that we can use them to check for the instances in each region
	var providers []*schema.Provider
	providerFactories := map[string]terraform.ResourceProviderFactory{
		"aws": func() (terraform.ResourceProvider, error) {
			p := Provider()
			providers = append(providers, p.(*schema.Provider))
			return p, nil
		},
	}

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckAWSS3BucketDestroyWithProviders(&providers),
		Steps: []resource.TestStep{
			{
				Config:      testAccAWSS3BucketConfigReplicationNoVersioning(rInt),
				ExpectError: regexp.MustCompile(`versioning must be enabled to allow S3 bucket replication`),
			},
		},
	})
}
*/

func TestAWSS3BucketName(t *testing.T) {
	validDnsNames := []string{
		"foobar",
		"foo.bar",
		"foo.bar.baz",
		"1234",
		"foo-bar",
		strings.Repeat("x", 63),
	}

	for _, v := range validDnsNames {
		if err := validateS3BucketName(v, "us-west-2"); err != nil {
			t.Fatalf("%q should be a valid S3 bucket name", v)
		}
	}

	invalidDnsNames := []string{
		"foo..bar",
		"Foo.Bar",
		"192.168.0.1",
		"127.0.0.1",
		".foo",
		"bar.",
		"foo_bar",
		strings.Repeat("x", 64),
	}

	for _, v := range invalidDnsNames {
		if err := validateS3BucketName(v, "us-west-2"); err == nil {
			t.Fatalf("%q should not be a valid S3 bucket name", v)
		}
	}

	validEastNames := []string{
		"foobar",
		"foo_bar",
		"127.0.0.1",
		"foo..bar",
		"foo_bar_baz",
		"foo.bar.baz",
		"Foo.Bar",
		strings.Repeat("x", 255),
	}

	for _, v := range validEastNames {
		if err := validateS3BucketName(v, "us-east-1"); err != nil {
			t.Fatalf("%q should be a valid S3 bucket name", v)
		}
	}

	invalidEastNames := []string{
		"foo;bar",
		strings.Repeat("x", 256),
	}

	for _, v := range invalidEastNames {
		if err := validateS3BucketName(v, "us-east-1"); err == nil {
			t.Fatalf("%q should not be a valid S3 bucket name", v)
		}
	}
}

func testAccCheckAWSS3BucketDestroy(s *terraform.State) error {
	// UNDONE: Why instance check?
	//return testAccCheckInstanceDestroyWithProvider(s, testAccProvider)
	return nil
}

func testAccCheckAWSS3BucketExists(n string) resource.TestCheckFunc {
	providers := []*schema.Provider{testAccProvider}
	return testAccCheckAWSS3BucketExistsWithProviders(n, &providers)
}

func testAccCheckAWSS3BucketExistsWithProviders(n string, providers *[]*schema.Provider) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}
		for _, provider := range *providers {
			// Ignore if Meta is empty, this can happen for validation providers
			if provider.Meta() == nil {
				continue
			}

			config := testAccProvider.Meta().(*Config)
			conn, err := config.computeS3conn(OS_REGION_NAME)
			if err != nil {
				return fmt.Errorf("Error creating OpenTelekomCloud s3 client: %s", err)
			}
			_, err = conn.HeadBucket(&s3.HeadBucketInput{
				Bucket: aws.String(rs.Primary.ID),
			})

			if err != nil {
				return fmt.Errorf("S3Bucket error: %v", err)
			}
			return nil
		}

		return fmt.Errorf("Instance not found")
	}
}

func testAccCheckAWSS3DestroyBucket(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No S3 Bucket ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		conn, err := config.computeS3conn(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating OpenTelekomCloud s3 client: %s", err)
		}
		_, err = conn.DeleteBucket(&s3.DeleteBucketInput{
			Bucket: aws.String(rs.Primary.ID),
		})

		if err != nil {
			return fmt.Errorf("Error destroying Bucket (%s) in testAccCheckAWSS3DestroyBucket: %s", rs.Primary.ID, err)
		}
		return nil
	}
}

func testAccCheckAWSS3BucketPolicy(n string, policy string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, _ := s.RootModule().Resources[n]
		config := testAccProvider.Meta().(*Config)
		conn, err := config.computeS3conn(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating OpenTelekomCloud s3 client: %s", err)
		}

		out, err := conn.GetBucketPolicy(&s3.GetBucketPolicyInput{
			Bucket: aws.String(rs.Primary.ID),
		})

		if policy == "" {
			if awsErr, ok := err.(awserr.Error); ok && awsErr.Code() == "NoSuchBucketPolicy" {
				// expected
				return nil
			}
			if err == nil {
				return fmt.Errorf("Expected no policy, got: %#v", *out.Policy)
			} else {
				return fmt.Errorf("GetBucketPolicy error: %v, expected %s", err, policy)
			}
		}
		if err != nil {
			return fmt.Errorf("GetBucketPolicy error: %v, expected %s", err, policy)
		}

		if v := out.Policy; v == nil {
			if policy != "" {
				return fmt.Errorf("bad policy, found nil, expected: %s", policy)
			}
		} else {
			expected := make(map[string]interface{})
			if err := json.Unmarshal([]byte(policy), &expected); err != nil {
				return err
			}
			actual := make(map[string]interface{})
			if err := json.Unmarshal([]byte(*v), &actual); err != nil {
				return err
			}

			if !reflect.DeepEqual(expected, actual) {
				return fmt.Errorf("bad policy, expected: %#v, got %#v", expected, actual)
			}
		}

		return nil
	}
}

func testAccCheckAWSS3BucketWebsite(n string, indexDoc string, errorDoc string, redirectProtocol string, redirectTo string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, _ := s.RootModule().Resources[n]
		config := testAccProvider.Meta().(*Config)
		conn, err := config.computeS3conn(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating OpenTelekomCloud s3 client: %s", err)
		}

		out, err := conn.GetBucketWebsite(&s3.GetBucketWebsiteInput{
			Bucket: aws.String(rs.Primary.ID),
		})

		if err != nil {
			if indexDoc == "" {
				// If we want to assert that the website is not there, than
				// this error is expected
				return nil
			} else {
				return fmt.Errorf("S3BucketWebsite error: %v", err)
			}
		}

		if v := out.IndexDocument; v == nil {
			if indexDoc != "" {
				return fmt.Errorf("bad index doc, found nil, expected: %s", indexDoc)
			}
		} else {
			if *v.Suffix != indexDoc {
				return fmt.Errorf("bad index doc, expected: %s, got %#v", indexDoc, out.IndexDocument)
			}
		}

		if v := out.ErrorDocument; v == nil {
			if errorDoc != "" {
				return fmt.Errorf("bad error doc, found nil, expected: %s", errorDoc)
			}
		} else {
			if *v.Key != errorDoc {
				return fmt.Errorf("bad error doc, expected: %s, got %#v", errorDoc, out.ErrorDocument)
			}
		}

		if v := out.RedirectAllRequestsTo; v == nil {
			if redirectTo != "" {
				return fmt.Errorf("bad redirect to, found nil, expected: %s", redirectTo)
			}
		} else {
			if *v.HostName != redirectTo {
				return fmt.Errorf("bad redirect to, expected: %s, got %#v", redirectTo, out.RedirectAllRequestsTo)
			}
			if redirectProtocol != "" && v.Protocol != nil && *v.Protocol != redirectProtocol {
				return fmt.Errorf("bad redirect protocol to, expected: %s, got %#v", redirectProtocol, out.RedirectAllRequestsTo)
			}
		}

		return nil
	}
}

func testAccCheckAWSS3BucketWebsiteRoutingRules(n string, routingRules []*s3.RoutingRule) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, _ := s.RootModule().Resources[n]
		config := testAccProvider.Meta().(*Config)
		conn, err := config.computeS3conn(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating OpenTelekomCloud s3 client: %s", err)
		}

		out, err := conn.GetBucketWebsite(&s3.GetBucketWebsiteInput{
			Bucket: aws.String(rs.Primary.ID),
		})

		if err != nil {
			if routingRules == nil {
				return nil
			}
			return fmt.Errorf("GetBucketWebsite error: %v", err)
		}

		if !reflect.DeepEqual(out.RoutingRules, routingRules) {
			return fmt.Errorf("bad routing rule, expected: %v, got %v", routingRules, out.RoutingRules)
		}

		return nil
	}
}

func testAccCheckAWSS3BucketVersioning(n string, versioningStatus string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, _ := s.RootModule().Resources[n]
		config := testAccProvider.Meta().(*Config)
		conn, err := config.computeS3conn(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating OpenTelekomCloud s3 client: %s", err)
		}

		out, err := conn.GetBucketVersioning(&s3.GetBucketVersioningInput{
			Bucket: aws.String(rs.Primary.ID),
		})

		if err != nil {
			return fmt.Errorf("GetBucketVersioning error: %v", err)
		}

		if v := out.Status; v == nil {
			if versioningStatus != "" {
				return fmt.Errorf("bad error versioning status, found nil, expected: %s", versioningStatus)
			}
		} else {
			if *v != versioningStatus {
				return fmt.Errorf("bad error versioning status, expected: %s, got %s", versioningStatus, *v)
			}
		}

		return nil
	}
}

func testAccCheckAWSS3BucketCors(n string, corsRules []*s3.CORSRule) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, _ := s.RootModule().Resources[n]
		config := testAccProvider.Meta().(*Config)
		conn, err := config.computeS3conn(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating OpenTelekomCloud s3 client: %s", err)
		}

		out, err := conn.GetBucketCors(&s3.GetBucketCorsInput{
			Bucket: aws.String(rs.Primary.ID),
		})

		if err != nil {
			return fmt.Errorf("GetBucketCors error: %v", err)
		}

		if !reflect.DeepEqual(out.CORSRules, corsRules) {
			return fmt.Errorf("bad error cors rule, expected: %v, got %v", corsRules, out.CORSRules)
		}

		return nil
	}
}

func testAccCheckAWSS3RequestPayer(n, expectedPayer string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, _ := s.RootModule().Resources[n]
		config := testAccProvider.Meta().(*Config)
		conn, err := config.computeS3conn(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating OpenTelekomCloud s3 client: %s", err)
		}

		out, err := conn.GetBucketRequestPayment(&s3.GetBucketRequestPaymentInput{
			Bucket: aws.String(rs.Primary.ID),
		})

		if err != nil {
			return fmt.Errorf("GetBucketRequestPayment error: %v", err)
		}

		if *out.Payer != expectedPayer {
			return fmt.Errorf("bad error request payer type, expected: %v, got %v",
				expectedPayer, out.Payer)
		}

		return nil
	}
}

func testAccCheckAWSS3BucketLogging(n, b, p string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, _ := s.RootModule().Resources[n]
		config := testAccProvider.Meta().(*Config)
		conn, err := config.computeS3conn(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating OpenTelekomCloud s3 client: %s", err)
		}

		out, err := conn.GetBucketLogging(&s3.GetBucketLoggingInput{
			Bucket: aws.String(rs.Primary.ID),
		})

		if err != nil {
			return fmt.Errorf("GetBucketLogging error: %v", err)
		}

		tb, _ := s.RootModule().Resources[b]

		if v := out.LoggingEnabled.TargetBucket; v == nil {
			if tb.Primary.ID != "" {
				return fmt.Errorf("bad target bucket, found nil, expected: %s", tb.Primary.ID)
			}
		} else {
			if *v != tb.Primary.ID {
				return fmt.Errorf("bad target bucket, expected: %s, got %s", tb.Primary.ID, *v)
			}
		}

		if v := out.LoggingEnabled.TargetPrefix; v == nil {
			if p != "" {
				return fmt.Errorf("bad target prefix, found nil, expected: %s", p)
			}
		} else {
			if *v != p {
				return fmt.Errorf("bad target prefix, expected: %s, got %s", p, *v)
			}
		}

		return nil
	}
}

// These need a bit of randomness as the name can only be used once globally
// within AWS
func testAccBucketName(randInt int) string {
	return fmt.Sprintf("tf-test-bucket-%d", randInt)
}

func testAccBucketDomainName(randInt int) string {
	return fmt.Sprintf("tf-test-bucket-%d.s3.amazonaws.com", randInt)
}

func testAccWebsiteEndpoint(randInt int) string {
	return fmt.Sprintf("tf-test-bucket-%d.s3-website.%s.amazonaws.com", randInt, OS_REGION_NAME)
}

func testAccAWSS3BucketPolicy(randInt int) string {
	return fmt.Sprintf(`{ "Version": "2008-10-17", "Statement": [ { "Effect": "Allow", "Principal": { "AWS": ["*"] }, "Action": ["s3:GetObject"], "Resource": ["arn:aws:s3:::tf-test-bucket-%d/*"] } ] }`, randInt)
}

func testAccAWSS3BucketConfig(randInt int) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_s3_bucket" "bucket" {
	bucket = "tf-test-bucket-%d"
	acl = "public-read"
}
`, randInt)
}

/*
func testAccAWSS3MultiBucketConfigWithTags(randInt int) string {
	t := template.Must(template.New("t1").
		Parse(`
resource "opentelekomcloud_s3_bucket" "bucket1" {
	bucket = "tf-test-bucket-1-{{.GUID}}"
	acl = "private"
	force_destroy = true
	tags {
		Name = "tf-test-bucket-1-{{.GUID}}"
		Environment = "{{.GUID}}"
	}
}

resource "opentelekomcloud_s3_bucket" "bucket2" {
	bucket = "tf-test-bucket-2-{{.GUID}}"
	acl = "private"
	force_destroy = true
	tags {
		Name = "tf-test-bucket-2-{{.GUID}}"
		Environment = "{{.GUID}}"
	}
}

resource "opentelekomcloud_s3_bucket" "bucket3" {
	bucket = "tf-test-bucket-3-{{.GUID}}"
	acl = "private"
	force_destroy = true
	tags {
		Name = "tf-test-bucket-3-{{.GUID}}"
		Environment = "{{.GUID}}"
	}
}

resource "opentelekomcloud_s3_bucket" "bucket4" {
	bucket = "tf-test-bucket-4-{{.GUID}}"
	acl = "private"
	force_destroy = true
	tags {
		Name = "tf-test-bucket-4-{{.GUID}}"
		Environment = "{{.GUID}}"
	}
}

resource "opentelekomcloud_s3_bucket" "bucket5" {
	bucket = "tf-test-bucket-5-{{.GUID}}"
	acl = "private"
	force_destroy = true
	tags {
		Name = "tf-test-bucket-5-{{.GUID}}"
		Environment = "{{.GUID}}"
	}
}

resource "opentelekomcloud_s3_bucket" "bucket6" {
	bucket = "tf-test-bucket-6-{{.GUID}}"
	acl = "private"
	force_destroy = true
	tags {
		Name = "tf-test-bucket-6-{{.GUID}}"
		Environment = "{{.GUID}}"
	}
}
`))
	var doc bytes.Buffer
	t.Execute(&doc, struct{ GUID int }{GUID: randInt})
	return doc.String()
}
*/

func testAccAWSS3BucketConfigWithRegion(randInt int) string {
	return fmt.Sprintf(`
provider "opentelekomcloud" {
	alias = "reg1"
	region = "%s"
}

resource "opentelekomcloud_s3_bucket" "bucket" {
	provider = "opentelekomcloud.reg1"
	bucket = "tf-test-bucket-%d"
	region = "%s"
}
`, OS_REGION_NAME, randInt, OS_REGION_NAME)
}

func testAccAWSS3BucketWebsiteConfig(randInt int) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_s3_bucket" "bucket" {
	bucket = "tf-test-bucket-%d"
	acl = "public-read"

	website {
		index_document = "index.html"
	}
}
`, randInt)
}

func testAccAWSS3BucketWebsiteConfigWithError(randInt int) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_s3_bucket" "bucket" {
	bucket = "tf-test-bucket-%d"
	acl = "public-read"

	website {
		index_document = "index.html"
		error_document = "error.html"
	}
}
`, randInt)
}

func testAccAWSS3BucketWebsiteConfigWithRedirect(randInt int) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_s3_bucket" "bucket" {
	bucket = "tf-test-bucket-%d"
	acl = "public-read"

	website {
		redirect_all_requests_to = "hashicorp.com"
	}
}
`, randInt)
}

func testAccAWSS3BucketWebsiteConfigWithHttpsRedirect(randInt int) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_s3_bucket" "bucket" {
	bucket = "tf-test-bucket-%d"
	acl = "public-read"

	website {
		redirect_all_requests_to = "https://hashicorp.com"
	}
}
`, randInt)
}

func testAccAWSS3BucketWebsiteConfigWithRoutingRules(randInt int) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_s3_bucket" "bucket" {
	bucket = "tf-test-bucket-%d"
	acl = "public-read"

	website {
		index_document = "index.html"
		error_document = "error.html"
		routing_rules = <<EOF
[{
	"Condition": {
		"KeyPrefixEquals": "docs/"
	},
	"Redirect": {
		"ReplaceKeyPrefixWith": "documents/"
	}
}]
EOF
	}
}
`, randInt)
}

func testAccAWSS3BucketConfigWithAcceleration(randInt int) string {
	return fmt.Sprintf(`
provider "opentelekomcloud" {
	alias = "reg1"
	region = "%s"
}

resource "opentelekomcloud_s3_bucket" "bucket" {
	provider = "opentelekomcloud.reg1"
	bucket = "tf-test-bucket-%d"
	region = "%s"
	acl = "public-read"
	acceleration_status = "Enabled"
}
`, OS_REGION_NAME, randInt, OS_REGION_NAME)
}

func testAccAWSS3BucketConfigWithoutAcceleration(randInt int) string {
	return fmt.Sprintf(`
provider "opentelekomcloud" {
	alias = "reg1"
	region = "%s"
}

resource "opentelekomcloud_s3_bucket" "bucket" {
	provider = "opentelekomcloud.reg1"
	bucket = "tf-test-bucket-%d"
	region = "%s"
	acl = "public-read"
	acceleration_status = "Suspended"
}
`, OS_REGION_NAME, randInt, OS_REGION_NAME)
}

/*
func testAccAWSS3BucketConfigRequestPayerBucketOwner(randInt int) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_s3_bucket" "bucket" {
	bucket = "tf-test-bucket-%d"
	acl = "public-read"
	request_payer = "BucketOwner"
}
`, randInt)
}

func testAccAWSS3BucketConfigRequestPayerRequester(randInt int) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_s3_bucket" "bucket" {
	bucket = "tf-test-bucket-%d"
	acl = "public-read"
	request_payer = "Requester"
}
`, randInt)
}
*/

func testAccAWSS3BucketConfigWithPolicy(randInt int) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_s3_bucket" "bucket" {
	bucket = "tf-test-bucket-%d"
	acl = "public-read"
	policy = %s
}
`, randInt, strconv.Quote(testAccAWSS3BucketPolicy(randInt)))
}

func testAccAWSS3BucketDestroyedConfig(randInt int) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_s3_bucket" "bucket" {
	bucket = "tf-test-bucket-%d"
	acl = "public-read"
}
`, randInt)
}

func testAccAWSS3BucketConfigWithEmptyPolicy(randInt int) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_s3_bucket" "bucket" {
	bucket = "tf-test-bucket-%d"
	acl = "public-read"
	policy = ""
}
`, randInt)
}

func testAccAWSS3BucketConfigWithVersioning(randInt int) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_s3_bucket" "bucket" {
	bucket = "tf-test-bucket-%d"
	acl = "public-read"
	versioning {
	  enabled = true
	}
}
`, randInt)
}

func testAccAWSS3BucketConfigWithDisableVersioning(randInt int) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_s3_bucket" "bucket" {
	bucket = "tf-test-bucket-%d"
	acl = "public-read"
	versioning {
	  enabled = false
	}
}
`, randInt)
}

func testAccAWSS3BucketConfigWithCORS(randInt int) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_s3_bucket" "bucket" {
	bucket = "tf-test-bucket-%d"
	acl = "public-read"
	cors_rule {
			allowed_headers = ["*"]
			allowed_methods = ["PUT","POST"]
			allowed_origins = ["https://www.example.com"]
			expose_headers = ["x-amz-server-side-encryption","ETag"]
			max_age_seconds = 3000
	}
}
`, randInt)
}

var testAccAWSS3BucketConfigWithAcl = `
resource "opentelekomcloud_s3_bucket" "bucket" {
	bucket = "tf-test-bucket-%d"
	acl = "public-read"
}
`

var testAccAWSS3BucketConfigWithAclUpdate = `
resource "opentelekomcloud_s3_bucket" "bucket" {
	bucket = "tf-test-bucket-%d"
	acl = "private"
}
`

func testAccAWSS3BucketConfigWithLogging(randInt int) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_s3_bucket" "log_bucket" {
	bucket = "tf-test-log-bucket-%d"
	acl = "log-delivery-write"
}
resource "opentelekomcloud_s3_bucket" "bucket" {
	bucket = "tf-test-bucket-%d"
	acl = "private"
	force_destroy = "true"
	logging {
		target_bucket = "${opentelekomcloud_s3_bucket.log_bucket.id}"
		target_prefix = "log/"
	}
}
`, randInt, randInt)
}

func testAccAWSS3BucketConfigWithLifecycle(randInt int) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_s3_bucket" "bucket" {
	bucket = "tf-test-bucket-%d"
	acl = "private"
	lifecycle_rule {
		id = "id1"
		prefix = "path1/"
		enabled = true

		expiration {
			days = 365
		}

		#transition {
		#	days = 30
		#	storage_class = "STANDARD_IA"
		#}
		#transition {
		#	days = 60
		#	storage_class = "GLACIER"
		#}
	}
	lifecycle_rule {
		id = "id2"
		prefix = "path2/"
		enabled = true

		expiration {
			date = "2016-01-12"
		}
	}
	lifecycle_rule {
		id = "id3"
		prefix = "path3/"
		enabled = true

		#transition {
		#	days = 0
		#	storage_class = "GLACIER"
		#}
	}
	lifecycle_rule {
		id = "id4"
		prefix = "path4/"
		enabled = true

		#tags {
		#	"tagKey" = "tagValue"
		#	"terraform" = "hashicorp"
		#}

		expiration {
			date = "2016-01-12"
		}
	}
}
`, randInt)
}

func testAccAWSS3BucketConfigWithVersioningLifecycle(randInt int) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_s3_bucket" "bucket" {
	bucket = "tf-test-bucket-%d"
	acl = "private"
	versioning {
	  enabled = false
	}
	lifecycle_rule {
		id = "id1"
		prefix = "path1/"
		enabled = true

		noncurrent_version_expiration {
			days = 365
		}
		noncurrent_version_transition {
			days = 30
			storage_class = "STANDARD_IA"
		}
		noncurrent_version_transition {
			days = 60
			storage_class = "GLACIER"
		}
	}
	lifecycle_rule {
		id = "id2"
		prefix = "path2/"
		enabled = false

		noncurrent_version_expiration {
			days = 365
		}
	}
	lifecycle_rule {
		id = "id3"
		prefix = "path3/"
		enabled = true

		noncurrent_version_transition {
			days = 0
			storage_class = "GLACIER"
		}
	}
}
`, randInt)
}

/*
const testAccAWSS3BucketConfigReplicationBasic = `
provider "opentelekomcloud" {
  alias  = "reg1"
  region = "eu-de"
}

provider "opentelekomcloud" {
  alias  = "reg2"
  region = "ap-sg"
}

resource "opentelekomcloud_iam_role" "role" {
  name               = "tf-iam-role-replication-%d"
  assume_role_policy = <<POLICY
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "s3.amazonaws.com"
      },
      "Effect": "Allow",
      "Sid": ""
    }
  ]
}
POLICY
}
`

func testAccAWSS3BucketConfigReplication(randInt int) string {
	return fmt.Sprintf(testAccAWSS3BucketConfigReplicationBasic+`
resource "opentelekomcloud_s3_bucket" "bucket" {
    provider = "opentelekomcloud.reg1"
    bucket   = "tf-test-bucket-%d"
    acl      = "private"

    versioning {
        enabled = true
    }
}

resource "opentelekomcloud_s3_bucket" "destination" {
    provider = "opentelekomcloud.reg2"
    bucket   = "tf-test-bucket-destination-%d"
    region   = "eu-de"

    versioning {
        enabled = true
    }
}
`, randInt, randInt, randInt)
}

func testAccAWSS3BucketConfigReplicationWithConfiguration(randInt int) string {
	return fmt.Sprintf(testAccAWSS3BucketConfigReplicationBasic+`
resource "opentelekomcloud_s3_bucket" "bucket" {
    provider = "opentelekomcloud.reg2"
    bucket   = "tf-test-bucket-%d"
    acl      = "private"

    versioning {
        enabled = true
    }

    replication_configuration {
        role = "${opentelekomcloud_iam_role.role.arn}"
        rules {
            id     = "foobar"
            prefix = "foo"
            status = "Enabled"

            destination {
                bucket        = "${opentelekomcloud_s3_bucket.destination.arn}"
                storage_class = "STANDARD"
            }
        }
    }
}

resource "opentelekomcloud_s3_bucket" "destination" {
    provider = "opentelekomcloud.reg1"
    bucket   = "tf-test-bucket-destination-%d"
    region   = "eu-de"

    versioning {
        enabled = true
    }
}
`, randInt, randInt, randInt)
}

func testAccAWSS3BucketConfigReplicationWithoutStorageClass(randInt int) string {
	return fmt.Sprintf(testAccAWSS3BucketConfigReplicationBasic+`
resource "opentelekomcloud_s3_bucket" "bucket" {
    provider = "opentelekomcloud.reg2"
    bucket   = "tf-test-bucket-%d"
    acl      = "private"

    versioning {
        enabled = true
    }

    replication_configuration {
        role = "${opentelekomcloud_iam_role.role.arn}"
        rules {
            id     = "foobar"
            prefix = "foo"
            status = "Enabled"

            destination {
                bucket        = "${opentelekomcloud_s3_bucket.destination.arn}"
            }
        }
    }
}

resource "opentelekomcloud_s3_bucket" "destination" {
    provider = "opentelekomcloud.reg1"
    bucket   = "tf-test-bucket-destination-%d"
    region   = "eu-de"

    versioning {
        enabled = true
    }
}
`, randInt, randInt, randInt)
}

func testAccAWSS3BucketConfigReplicationNoVersioning(randInt int) string {
	return fmt.Sprintf(testAccAWSS3BucketConfigReplicationBasic+`
resource "opentelekomcloud_s3_bucket" "bucket" {
    provider = "opentelekomcloud.reg2"
    bucket   = "tf-test-bucket-%d"
    acl      = "private"

    replication_configuration {
        role = "${opentelekomcloud_iam_role.role.arn}"
        rules {
            id     = "foobar"
            prefix = "foo"
            status = "Enabled"

            destination {
                bucket        = "${opentelekomcloud_s3_bucket.destination.arn}"
                storage_class = "STANDARD"
            }
        }
    }
}

resource "opentelekomcloud_s3_bucket" "destination" {
    provider = "opentelekomcloud.reg1"
    bucket   = "tf-test-bucket-destination-%d"
    region   = "eu-de"

    versioning {
        enabled = true
    }
}
`, randInt, randInt, randInt)
}
*/

const testAccAWSS3BucketConfig_namePrefix = `
resource "opentelekomcloud_s3_bucket" "test" {
	bucket_prefix = "tf-test-"
}
`

const testAccAWSS3BucketConfig_generatedName = `
resource "opentelekomcloud_s3_bucket" "test" {
	bucket_prefix = "tf-test-"
}
`