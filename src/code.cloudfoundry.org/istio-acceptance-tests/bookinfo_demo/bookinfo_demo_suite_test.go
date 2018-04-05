package bookinfo_demo

import (
	"fmt"
	"os"
	"testing"
	"time"

	"code.cloudfoundry.org/istio-acceptance-tests/config"
	"github.com/cloudfoundry-incubator/cf-test-helpers/cf"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gexec"
	"github.com/sclevine/agouti"
)

func TestBookinfoDemo(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "BookinfoDemo Suite")
}

var (
	agoutiDriver   *agouti.WebDriver
	c              config.Config
	defaultTimeout = 360 * time.Second
	org, space     = "ISTO-ORG", "ISTIO-SPACE"
)

var _ = BeforeSuite(func() {
	var err error
	configPath := os.Getenv("CONFIG")
	Expect(configPath).NotTo(BeEmpty())
	c, err = config.NewConfig(configPath)
	Expect(err).ToNot(HaveOccurred())

	enableDockerCmd := cf.Cf("enable-feature-flag", "diego_docker").Wait(defaultTimeout)
	Expect(enableDockerCmd).To(Exit(0))

	apiCmd := cf.Cf("api", fmt.Sprintf("api.%s", c.ApiEndpoint), "--skip-ssl-validation").Wait(defaultTimeout)
	Expect(apiCmd).To(Exit(0))
	loginCmd := cf.Cf("auth", c.AdminUser, c.AdminPassword).Wait(defaultTimeout)
	Expect(loginCmd).To(Exit(0))

	orgCmd := cf.Cf("create-org", org).Wait(defaultTimeout)
	Expect(orgCmd).To(Exit(0))
	spaceCmd := cf.Cf("create-space", "-o", org, space).Wait(defaultTimeout)
	Expect(spaceCmd).To(Exit(0))
	targetCmd := cf.Cf("target", "-o", org, "-s", space).Wait(defaultTimeout)
	Expect(targetCmd).To(Exit(0))

	//TODO: fix the dockerhub to cf-routing and also the image tags
	productPagePush := cf.Cf("push", "productpage", "-o", "zlav/examples-bookinfo-productpage-v1:1.0.0", "-d", "bosh-lite.com").Wait(defaultTimeout)
	Expect(productPagePush).To(Exit(0))
	ratingsPush := cf.Cf("push", "ratings", "-o", "zlav/examples-bookinfo-ratings-v1:1.0.0", "-d", "apps.internal").Wait(defaultTimeout)
	Expect(ratingsPush).To(Exit(0))
	reviewsPush := cf.Cf("push", "ratings", "-o", "zlav/examples-bookinfo-reviews-v3:1.0.0", "-d", "apps.internal").Wait(defaultTimeout)
	Expect(reviewsPush).To(Exit(0))
	detailsPush := cf.Cf("push", "details", "-o", "zlav/examples-bookinfo-details-v1:1.0.0", "-d", "apps.internal").Wait(defaultTimeout)
	Expect(detailsPush).To(Exit(0))

	productDetailsPolicy := cf.Cf("add-network-policy", "productpage", "--destination-app", "details", "--protocol", "tcp", "--port", "9080").Wait(defaultTimeout)
	Expect(productDetailsPolicy).To(Exit(0))
	productReviewsPolicy := cf.Cf("add-network-policy", "productpage", "--destination-app", "reviews", "--protocol", "tcp", "--port", "9080").Wait(defaultTimeout)
	Expect(productReviewsPolicy).To(Exit(0))
	reviewsRatingsPolicy := cf.Cf("add-network-policy", "reviews", "--destination-app", "ratings", "--protocol", "tcp", "--port", "9080").Wait(defaultTimeout)
	Expect(reviewsRatingsPolicy).To(Exit(0))

	agoutiDriver = agouti.ChromeDriver(
		agouti.ChromeOptions("args", []string{
			"--headless",
			"--disable-gpu",
			"--allow-insecure-localhost",
		}),
	)

	Expect(agoutiDriver.Start()).To(Succeed())
})

var _ = AfterSuite(func() {
	cleanUpCmd := cf.Cf("delete-org", org, "-f").Wait(defaultTimeout)
	Expect(cleanUpCmd).To(Exit(0))
	Expect(agoutiDriver.Stop()).To(Succeed())
})
