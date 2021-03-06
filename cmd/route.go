/*
Copyright © 2020 Mohamed Beydoun mohamed.beydoun@mail.mcgill.ca

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/MohamedBeydoun/atlas/pkg/generator"
	"github.com/iancoleman/strcase"
	"github.com/spf13/cobra"
)

// routeCmd represents the route command
var routeCmd = &cobra.Command{
	Use:   "route [flags]",
	Short: "Generates an express route.",
	Long:  `Generates a new express route with it's respective controller functions.`,
	RunE:  generateRoute,
}

func init() {
	generateCmd.AddCommand(routeCmd)
	routeCmd.Flags().StringP("router", "r", "", "Router name")
	routeCmd.Flags().StringP("method", "m", "get", "HTTP method for the route")
	routeCmd.Flags().StringP("url", "u", "/dummy", "Route endpoint")
	routeCmd.Flags().StringP("controller", "c", "dummy", "The controller function name for handling the route logic")
	routeCmd.MarkFlagRequired("router")
}

func generateRoute(cmd *cobra.Command, args []string) error {
	router, err := cmd.Flags().GetString("router")
	if err != nil {
		return err
	}
	router = strcase.ToLowerCamel(router)

	method, err := cmd.Flags().GetString("method")
	if err != nil {
		return err
	}
	method = strings.ToLower(method)

	url, err := cmd.Flags().GetString("url")
	if err != nil {
		return err
	}

	controller, err := cmd.Flags().GetString("controller")
	if err != nil {
		return err
	}
	controller = strcase.ToLowerCamel(controller)

	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	// check if router exists
	if _, err := os.Stat(fmt.Sprintf("%s/src/routes/%s.ts", wd, router)); os.IsNotExist(err) {
		return fmt.Errorf("router %s does not exist", router)
	}

	// check if valid http method
	expectedHTTPMethods := []string{"get", "post", "put", "patch", "update", "delete"}
	for _, expectedMethod := range expectedHTTPMethods {
		if method == expectedMethod {
			break
		}
		if method != expectedMethod && expectedMethod == "delete" {
			return fmt.Errorf("unknown HTTP method %s", method)
		}
	}

	// check if valid url
	validURL, err := regexp.MatchString(`^\/[/.a-zA-Z0-9-:]+$`, url)
	if err != nil {
		return err
	}
	if string(url[0]) != "/" || !validURL {
		return fmt.Errorf("invalid route format: %s", url)
	}

	route, err := generator.NewRoute(router, method, url, controller)
	if err != nil {
		return err
	}

	err = route.Create()
	if err != nil {
		return err
	}

	return nil
}
