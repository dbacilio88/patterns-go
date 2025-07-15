package app

import (
	"fmt"
	"github.com/dbacilio88/patterns-go/internal/config/app"
	"github.com/gorilla/mux"
	"log"
)

/**
 * startup
 * <p>
 * This file contains core data structures and logic used throughout the application.
 *
 * <p><strong>Copyright © 2025 – All rights reserved.</strong></p>
 *
 * <p>This source code is distributed under a collaborative license.</p>
 *
 * <p>
 * Contributions, suggestions, and improvements are welcome!
 * You are free to fork, modify, and submit pull requests under the terms of the repository's license.
 * Please ensure proper attribution to the original author(s) and preserve this notice in derivative works.
 * </p>
 *
 * @author Christian Bacilio De La Cruz
 * @email dbacilio88@outlook.es
 * @since 6/29/2025
 */

func Startup(configPath string) {
	log.Println("Startup")
	var err error
	app.Once.Do(func() {
		if err = ConfigureApplication(configPath); err != nil {
			message = fmt.Sprintf("Configuration Error: %s", err.Error())
			return
		}

		if err = ExecuteRabbitProcess(); err != nil {
			message = fmt.Sprintf("Rabbit Process Error: %s", err.Error())
			return
		}
	})
}

func setupHTTPServer(router *mux.Router) {
	router.Use()
}
