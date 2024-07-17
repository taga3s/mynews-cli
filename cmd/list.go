/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
	"fmt"
	"mynews/api"
	"mynews/model"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Fetch and display the latest news articles",
	Long:  `Fetches and displays the latest news articles.`,
	Run: func(cmd *cobra.Command, args []string) {
		run()
	},
}

func run() {
	newsResponse, err := api.Fetch()
	if err != nil {
		fmt.Println(err)
		return
	}
	newsList := make([]model.News, 0, len(newsResponse.Articles))

	for _, v := range newsResponse.Articles {
		news := model.News{
			Author:      v.Author,
			Title:       v.Title,
			Description: v.Description,
			PublishedAt: v.PublishedAt,
			URL:         v.URL,
		}
		newsList = append(newsList, news)
	}

	selectedNewsList, err := multiSelectNews(newsList)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = api.Notify(selectedNewsList)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Notify to your line app!")
}

func multiSelectNews(newsList []model.News) ([]model.News, error) {
	options := make([]huh.Option[model.News], 0, len(newsList))
	for _, v := range newsList {
		options = append(options, huh.NewOption(v.Title, v))
	}

	selectedNewsList := []model.News{}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewMultiSelect[model.News]().
				Options(options...).
				Title("News").
				Value(&selectedNewsList).
				Validate(validateMultiSelect),
		),
	)

	err := form.Run()
	if err != nil {
		return []model.News{}, err
	}

	return selectedNewsList, nil
}

func validateMultiSelect(selectedNewsList []model.News) error {
	if len(selectedNewsList) == 0 {
		return errors.New("You should select at least 1 article.")
	}
	return nil
}

func init() {
	rootCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
