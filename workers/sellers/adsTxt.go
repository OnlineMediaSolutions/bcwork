package sellers

import (
	"fmt"
	"github.com/m6yf/bcwork/utils/constant"
	"io"
	"net/http"
	"strings"
)

func (worker *Worker) prepareAdsTxtData(publisherDomains []PublisherDomain) []AdsTxt {
	adsTxtData := make([]AdsTxt, 0)
	jobs := make(chan PublisherDomain, len(publisherDomains))
	results := make(chan AdsTxt, len(publisherDomains))

	for i := 0; i < constant.SellersJsonWorkerCount; i++ {
		go worker.adsTxtWorkerWithStatus(jobs, results)
	}

	for _, data := range publisherDomains {
		jobs <- data
	}
	close(jobs)

	for i := 0; i < len(publisherDomains); i++ {
		result := <-results
		adsTxtData = append(adsTxtData, result)
	}
	close(results)

	return adsTxtData
}

func (worker *Worker) adsTxtWorkerWithStatus(jobs <-chan PublisherDomain, results chan<- AdsTxt) {
	for data := range jobs {
		status, err := worker.getAdsTxtStatus(data.Domain, data.Publisher, data.SellerId)
		if err != nil {
			fmt.Printf("Error validating domain %s: %v\n", data.Domain, err)
			results <- AdsTxt{
				Domain:        data.Domain,
				SellerId:      data.SellerId,
				PublisherName: data.Publisher,
				AdsTxtStatus:  "not verified",
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

func (worker *Worker) getAdsTxtStatus(domain, publisherName, sellerId string) (string, error) {
	url := fmt.Sprintf("https://%s/ads.txt", domain)

	resp, err := http.Get(url)
	if err != nil {
		return "not verified", fmt.Errorf("failed to fetch ads.txt for domain %s: %v", domain, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "not valid", fmt.Errorf("ads.txt not found or invalid for domain %s", domain)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "not verified", fmt.Errorf("failed to read ads.txt for domain %s: %v", domain, err)
	}

	content := string(body)
	publisherExists := strings.Contains(content, publisherName)
	sellerIdExists := strings.Contains(content, sellerId)

	if publisherExists && sellerIdExists {
		return "included", nil
	}
	return "not included", nil
}
