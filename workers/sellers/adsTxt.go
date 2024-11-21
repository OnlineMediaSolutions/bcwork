package sellers

import (
	"fmt"
	"github.com/m6yf/bcwork/utils/constant"
	"io"
	"net/http"
	"strings"
)

func (worker *Worker) adsTxtWorkerWithStatus(jobs <-chan PublisherDomain, results chan<- AdsTxt) {
	for data := range jobs {

		status, err := worker.getAdsTxtStatus(data.Domain, data.SellerId)
		if err != nil {
			results <- AdsTxt{
				Domain:        data.Domain,
				SellerId:      data.SellerId,
				PublisherName: data.Publisher,
				AdsTxtStatus:  constant.AdsTxtNotVerifiedStatus,
			}
			continue
		}

		results <- AdsTxt{
			Domain:        data.Domain,
			SellerId:      data.SellerId,
			PublisherName: data.Publisher,
			AdsTxtStatus:  status,
		}
	}
}

func (worker *Worker) getAdsTxtStatus(domain, sellerId string) (string, error) {
	url := fmt.Sprintf("https://%s/ads.txt", domain)

	resp, err := http.Get(url)
	if err != nil {
		return constant.AdsTxtNotVerifiedStatus, fmt.Errorf("failed to fetch ads.txt for domain %s: %v", domain, err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return constant.AdsTxtNotVerifiedStatus, fmt.Errorf("ads.txt not found or invalid for domain %s", domain)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return constant.AdsTxtNotVerifiedStatus, fmt.Errorf("failed to read ads.txt for domain %s: %v", domain, err)
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
			key := fmt.Sprintf("%s", currentSellerId)
			adsTxtMap[key] = struct{}{}
		}
	}

	searchKey := fmt.Sprintf("%s", strings.ToLower(strings.TrimSpace(sellerId)))
	if _, exists := adsTxtMap[searchKey]; exists {
		return constant.AdsTxtIncludedStatus, nil
	}

	return constant.AdsTxtNotIncludedStatus, nil

}

func (worker *Worker) enhancePublisherDomains(domains []PublisherDomain) []PublisherDomain {
	results := make([]PublisherDomain, 0, len(domains))
	jobs := make(chan PublisherDomain, len(domains))
	output := make(chan PublisherDomain, len(domains))

	numWorkers := constant.SellersJsonWorkerCount
	for i := 0; i < numWorkers; i++ {
		go func() {
			for domain := range jobs {
				status, err := worker.getAdsTxtStatus(domain.Domain, domain.SellerId)
				if err != nil {
					status = constant.AdsTxtNotVerifiedStatus
				}
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
