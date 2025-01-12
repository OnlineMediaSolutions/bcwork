package sellers

import (
	"fmt"
	"github.com/m6yf/bcwork/utils/constant"
	"io"
	"net/http"
	"strings"
	"time"
)

func (worker *Worker) GetAdsTxtStatus(domain, sellerId, competitorType string) string {

	if domain == "" {
		return constant.AdsTxtNotVerifiedStatus
	}

	url := worker.GetAdsTxtUrl(domain, competitorType)

	client := &http.Client{
		Timeout: constant.AdsTxtRequestTimeout * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		return constant.AdsTxtNotVerifiedStatus
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return constant.AdsTxtNotVerifiedStatus
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return constant.AdsTxtNotVerifiedStatus
	}

	content := strings.ToLower(string(body))
	adsTxtMap := make(map[string]struct{})

	lines := strings.Split(content, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		fields := strings.Split(line, ",")
		if len(fields) >= 2 {
			currentSellerId := strings.TrimSpace(fields[1])
			adsTxtMap[currentSellerId] = struct{}{}
		}
	}

	searchKey := strings.ToLower(strings.TrimSpace(sellerId))
	if _, exists := adsTxtMap[searchKey]; exists {
		return constant.AdsTxtIncludedStatus
	}

	return constant.AdsTxtNotIncludedStatus
}

func (worker *Worker) GetAdsTxtUrl(domain string, competitorType string) string {
	url := fmt.Sprintf("https://%s/ads.txt", domain)
	if competitorType == "inapp" {
		url = fmt.Sprintf("https://%s/app-ads.txt", domain)
	}
	return url
}

func (worker *Worker) enhancePublisherDomains(domains []PublisherDomain, competitorType string) []PublisherDomain {
	results := make([]PublisherDomain, 0, len(domains))
	jobs := make(chan PublisherDomain, len(domains))
	output := make(chan PublisherDomain, len(domains))

	numWorkers := constant.SellersJsonWorkerCount
	for i := 0; i < numWorkers; i++ {
		go func() {
			for domain := range jobs {
				status := worker.GetAdsTxtStatus(domain.Domain, domain.SellerId, competitorType)
				domain.AdsTxtStatus = status
				output <- domain
			}
		}()
	}

	for _, domain := range domains {
		jobs <- domain
	}
	close(jobs)

	for range domains {
		results = append(results, <-output)
	}
	close(output)

	return results
}
