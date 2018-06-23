#!/usr/bin/env Rscript
Sys.setlocale("LC_MESSAGES", 'en_GB.UTF-8')
Sys.setenv(LANG = "en_US.UTF-8")

options(repos = c(CRAN = "http://cran.rstudio.com"))
install.packages(c("devtools",
                   "ggplot2",
                   "dplyr",
                   "httr",
                   "stringr",
                   "jsonlite",
                   "magrittr",
                   "testthat",
                   "xlsx",
                   "knitr"),
                   dependencies = TRUE)