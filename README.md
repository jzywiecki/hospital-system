# Hospital Examination Management System
This project implements a hospital examination management system in Go, utilizing RabbitMQ for efficient communication between doctors and technicians. The system allows doctors to send examination requests (such as hip, knee, or elbow examinations) to technicians, who can then process the examinations and return the results. The system is designed to handle message exchanges through RabbitMQ’s topic-based routing, ensuring that examination requests and results are properly routed between medical staff.

Additionally, the system includes a logging mechanism that captures and broadcasts logs through a fanout exchange, providing real-time insights into system activities.

# System design with RabbitMQ

![schemat](https://github.com/user-attachments/assets/ee823343-44f1-4a06-a14c-b888957b1361)
