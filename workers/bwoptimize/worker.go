package bwoptimize

import (
	"context"
	"fmt"
	"github.com/digitalocean/godo"
	"github.com/m6yf/bcwork/config"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"strconv"
	"strings"
)

type Worker struct {
	Name     string `json:"name"`
	DoClient *godo.Client
}

func (w *Worker) Init(ctx context.Context, conf config.StringMap) error {

	w.DoClient = godo.NewFromToken(viper.GetString("digitalocean.token"))

	return nil
}

func (w *Worker) Do(ctx context.Context) error {

	return w.TerminateDroplets(ctx)

	//// Configuration for the droplets
	//region := "nyc3"                                           // Choose a region, e.g., nyc3
	//size := "s-2vcpu-4gb"                                      // $24/month droplet
	//image := godo.DropletCreateImage{Slug: "ubuntu-20-04-x64"} // Base image
	//sshKeys := []godo.DropletCreateSSHKey{}
	//
	//// Create 100 droplets
	//for i := 101; i <= 500; i++ {
	//	dropletName := fmt.Sprintf("bwoptimize-%03d", i)
	//
	//	// Define droplet create request
	//	createRequest := &godo.DropletCreateRequest{
	//		Name:    dropletName,
	//		Region:  region,
	//		Size:    size,
	//		Image:   image,
	//		SSHKeys: sshKeys,
	//		Tags:    []string{"bwoptimize"},
	//	}
	//
	//	// Create the droplet
	//	fmt.Printf("Creating droplet: %s\n", dropletName)
	//	droplet, _, err := w.DoClient.Droplets.Create(context.Background(), createRequest)
	//	if err != nil {
	//		log.Fatal().Msgf("Error creating droplet %s: %v", dropletName, err)
	//	}
	//
	//	// Print droplet details
	//	fmt.Printf("Created droplet %s with ID %d\n", droplet.Name, droplet.ID)
	//
	//	// Pause to avoid API rate limits
	//	time.Sleep(2 * time.Second)
	//}
	return nil
}

func (w *Worker) GetSleep() int {
	return 0
}

func (w *Worker) TerminateDroplets(ctx context.Context) error {

	for i := 1; i < 5; i++ {
		// Get all droplets
		droplets, _, err := w.DoClient.Droplets.List(context.Background(), &godo.ListOptions{
			Page:    i,
			PerPage: 200, // Adjust if you have more droplets

		})
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to list droplets")
		}

		// Filter droplets that match our pattern
		var targetDroplets []godo.Droplet
		// Iterate through droplets and delete those within the range
		for _, droplet := range droplets {
			name := droplet.Name

			// Check if name starts with "bwoptimize-"
			if strings.HasPrefix(name, "bwoptimize-") {
				// Extract the numeric suffix (e.g., "bwoptimize-250" → 250)
				parts := strings.Split(name, "-")
				if len(parts) < 2 {
					continue
				}

				// Convert suffix to an integer
				num, err := strconv.Atoi(parts[1])
				if err != nil {
					continue
				}

				// Check if the number falls within the range 250-500
				if num >= 250 && num <= 500 {
					fmt.Printf("Deleting droplet: %s (ID: %d)\n", name, droplet.ID)
					_, err := w.DoClient.Droplets.Delete(ctx, droplet.ID)
					if err != nil {
						fmt.Printf("Failed to delete droplet %s: %v\n", name, err)
					} else {
						fmt.Printf("Deleted droplet: %s\n", name)
					}
				}
			}
		}

		fmt.Printf("Found %d droplets to terminate\n", len(targetDroplets))

	}

	return nil
}

//
//func DropletList(ctx context.Context, client *godo.Client) ([]godo.Droplet, error) {
//	// create a list to hold our droplets
//	list := []godo.Droplet{}
//
//	// create options. initially, these will be blank
//	opt := &godo.ListOptions{}
//	for {
//		droplets, resp, err := client.Droplets.List(ctx, opt)
//		if err != nil {
//			return nil, err
//		}
//
//		// append the current page's droplets to our list
//		list = append(list, droplets...)
//
//		// if we are at the last page, break out the for loop
//		if resp.Links == nil || resp.Links.IsLastPage() {
//			break
//		}
//
//		page, err := resp.Links.CurrentPage()
//		if err != nil {
//			return nil, err
//		}
//
//		// set the page we want for the next request
//		opt.Page = page + 1
//	}
//
//	return list, nil
//}
