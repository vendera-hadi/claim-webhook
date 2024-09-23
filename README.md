# Claim Webhook Listener

## Overview

Claim Webhook Listener is an example of specialized webhook designed for insurance companies and healthcare facilities. It listens for requests from the Satusehat platform and prints the received messages and notifications to an executable file.

## Features

- **Real-time Listening**: Continuously listens for incoming requests from Satusehat.
- **Message Logging**: Logs all received messages and notifications.
- **Executable Output**: Prints messages and notifications to an executable file for further processing.

## Installation

1. **Clone the repository**:
    ```sh
    git clone https://github.com/vendera-hadi/claim-webhook.git
    cd claim-webhook
    ```

2. **Configure the environment**:
    - Create a `.env` file in the root directory.
    - Add your configuration settings (e.g., API keys, webhook URLs).
      PORT, PUBLIC_KEY_SS (base64), PRIVATE_KEY_ORG (base64)


## Usage

1. **Start the webhook listener**:
    ```sh
    PORT=1234 go run main.go
    ```
    or just build the executable file
    ```sh
    go build -o webhook.exe
    ```


2. **Monitor the logs**:
    - Check the console output for real-time logs.
    - Review the executable file for detailed message logs.

## Contributing

We welcome contributions! Please follow these steps:

1. Fork the repository.
2. Create a new branch (`git checkout -b feature-branch`).
3. Make your changes.
4. Commit your changes (`git commit -m 'Add new feature'`).
5. Push to the branch (`git push origin feature-branch`).
6. Open a pull request.

## License

This project is licensed under the MIT License. See the LICENSE file for details.

## Contact

For any questions or support, please open an issue or contact us at vendera.hadi@dto.kemkes.go.id
