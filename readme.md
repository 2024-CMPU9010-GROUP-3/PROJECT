<div align="center">

<img src="https://github.com/user-attachments/assets/c147b766-d1bd-4cf7-b4e9-fa49705c89b1" align="center" width="144px" height="144px"/>

# Magpie

_Repository for Group 3 of the TU Dublin ASD/DS masters Group Project_

</div>

<div align="center">

![GitHub Repo stars](https://img.shields.io/github/stars/2024-CMPU9010-GROUP-3/magpie?style=for-the-badge)
![GitHub forks](https://img.shields.io/github/forks/2024-CMPU9010-GROUP-3/magpie?style=for-the-badge)

</div>

# ❓ What is this repository for?

_Magpie_ is a geographical information service that allows Civil Planners and other users to easily explore public amenities 'at a glance'.

You can access the live version of the project [here](https://magpie.solonsstuff.com/).

![magpie](https://github.com/user-attachments/assets/bcffd0ca-e228-484c-9236-d749e9769932)

# 📂 Repository Structure

The project is divided into several parts:

1. `Backend`: Public and Private API services

2. `Distribution`: Docker Compose and Kubernetes deployment files

3. `Documentation`: Project documentation, including interim reports, MOMs, and retrospectives

4. `Frontend`: Next.js frontend

5. `Python`: Python scripts for data processing, machine learning, and data visualization

The file structure is as follows:

```sh
📁
├──📁 Backend
│   ├──📁 cmd
│   ├──📁 internal
│   └──📁 sql
├──📁 Distribution
│   ├──📁 compose
│   └──📁 kubernetes
├──📁 Documentation
│   ├──📁 gantt-chart
│   ├──📁 interim-report
│   ├──📁 mom
│   ├──📁 presentation
│   ├──📁 project plan
│   ├──📁 retrospectives
│   ├──📁 survey
│   └──📁 ux-documents
├──📁 Frontend
│   ├──📁 public
│   └──📁 src
└──📁 Python
    ├──📁 charts
    ├──📁 notebooks
    ├──📁 script
    └──📁 yolo
```

# ⚙️ Using the development containers

This repository is set up to use dev containers for development. To use them, you will need:

- to have [docker](https://www.docker.com/) installed on your machine
- to have the [dev container extension](https://marketplace.visualstudio.com/items?itemName=ms-vscode-remote.remote-containers) installed in VS Code

Once installed, open the command pallet with `Ctrl+Shift+P` and type `Dev Containers: Open Folder in Container`. Then select the desired folder.
