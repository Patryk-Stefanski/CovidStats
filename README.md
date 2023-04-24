# Covid-19 Statistics Web Application
This repository contains a simple web application written in Go that displays COVID-19 statistics for each country on a specific date. The application is designed to demonstrate a GitOps implementation as part of a final year project.

## Prerequisite
Api Key for Covid19 API from RapidAPI.com

## Installation
To install and run the application, please follow these steps:

1. Clone the repository using the command: git clone git@github.com:Patryk-Stefanski/CovidStats.git
2. Navigate to the root directory of the project: cd CovidStats
3. Build the Docker image using the command: podman build -t docker.io/<username>/covidstats:v.1.0.0 .
4. Run the Docker container using the command: podman run -i -p 3000:3000 -e COVID_STATS_API_KEY=<covid_api_key> covidstats:v1.0.0
5. Once the container is up and running, you can access the application by navigating to http://localhost:3000 in your web browser.

## Usage
The application provides a simple user interface that allows you to select a date and a country to display COVID-19 statistics. To use the application, follow these steps:

1. Select a date from the date picker.
2. Select a country from the dropdown list.
3. Click the "Show Statistics" button.
4. The application will display the confirmed cases, deaths, and recoveries for the selected country on the selected date.

## GitOps Implementation
As part of a final year project, this application was used to demonstrate a GitOps implementation. GitOps is a set of practices that use Git as a single source of truth for declarative infrastructure and applications. In this implementation, changes to the application code or configuration files trigger a series of automated processes that update the application in a Kubernetes cluster.
