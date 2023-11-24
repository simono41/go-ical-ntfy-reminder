# GO ical ntfy reminder

This is a simple Go application that does iCalendar files (.ics) analyzes and sends notifications by ntfy based on the events in these iCalendar files.

## Prerequisites

Before you begin, ensure you have the following installed:

- [Docker](https://docs.docker.com/get-docker/)
- [Docker Compose](https://docs.docker.com/compose/install/)

## Getting Started

1. **Clone the repository:**

    ```bash
    git clone https://code.brothertec.eu/simono41/go-ical-ntfy-reminder.git
    ```

2. **Create an environment variables file:**

   Create a file named `.env` in your project directory. Fill this file with the required environment variables:

    ```env
    TO_EMAIL=recipient@example.com (optional)
    NTFY_AUTH=echo "Basic $(echo -n 'testuser:fakepassword' | base64)"
    NTFY_HOST=https://ntfy.sh/alerts
    LOCATION=Europe/Berlin (default)
    ICS_DIR=/path/to/ics/files
    ```

   Customize the values according to your application.

3. **Build and run the Docker container:**

    ```bash
    docker compose up -d --build
    ```

   This command uses Docker Compose to build and run the container, loading environment variables from the `.env` file.

4. **Access your application:**

   Open your web browser and go to [http://localhost:8080](http://localhost:8080).

## Customizing the Docker Image

- If your application uses additional environment variables, add them to the `.env` file.
- Customize the Dockerfile or docker-compose.yml if needed.

## Running as a Cronjob

To run your application as a daily cronjob at 6 AM, add the following cronjob entry:

```cron
0 6 * * * docker compose -f /opt/containers/mail-reminder/docker-compose.yml up --build --exit-code-from go-app
```

This cronjob will execute the Docker Compose command daily at 6 AM.

## Contributing

If you'd like to contribute, please fork the repository and create a pull request. Feel free to open an issue if you encounter any problems or have suggestions.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
