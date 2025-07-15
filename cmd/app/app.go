package app

import (
	"flag"
	"fmt"
	"time"
)

/**
 * app
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

var configPath string

func Run() {
	start := time.Now()
	flag.StringVar(&configPath, "path", "../config", "config file path")
	fmt.Println("Hello World", configPath)
	flag.Parse()

	Startup(configPath)

	elapsed := time.Since(start)
	fmt.Println("Start Time: ", elapsed)
	select {
	case <-time.After(30 * time.Second):
	}
}
