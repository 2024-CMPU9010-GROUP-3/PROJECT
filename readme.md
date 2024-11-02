<div align="center">

<img src="https://github.com/user-attachments/assets/c147b766-d1bd-4cf7-b4e9-fa49705c89b1" align="center" width="144px" height="144px"/>

# Magpie: _Services at a Glance_

_Repository for Group 3 of the TU Dublin ASD/DS masters Group Project_

</div>

<div align="center">

![GitHub Actions Workflow Status](https://img.shields.io/github/actions/workflow/status/2024-CMPU9010-GROUP-3/magpie/frontend-build.yaml?style=flat-square&label=Frontend-Build)
![GitHub Actions Workflow Status](https://img.shields.io/github/actions/workflow/status/2024-CMPU9010-GROUP-3/magpie/backend-public-build.yaml?style=flat-square&label=Backend-Public-Build)
![GitHub Actions Workflow Status](https://img.shields.io/github/actions/workflow/status/2024-CMPU9010-GROUP-3/magpie/backend-private-build.yaml?style=flat-square&label=Backend-Private-Build)


![GitHub Repo stars](https://img.shields.io/github/stars/2024-CMPU9010-GROUP-3/magpie?style=for-the-badge)
![GitHub forks](https://img.shields.io/github/forks/2024-CMPU9010-GROUP-3/magpie?style=for-the-badge)

</div>

## ❓ What is this repository for?

_Magpie_ is a geographical information service that allows Civil Planners and other users to easily explore public amenities 'at a glance'.

<p align="center">
  <img src="https://github.com/user-attachments/assets/bcffd0ca-e228-484c-9236-d749e9769932" width="800"/>
</p>

You can access a hosted version of the project [here](https://magpie.solonsstuff.com/).

## ✨ Features

- **Geographical Information Service**: Easily explore public amenities such as parking spaces, bike sheds, and accessible ramps.
- **User Authentication**: Secure login and registration for users.
- **Data Integration**: Integrates various data sources for comprehensive information.
- **Interactive Maps**: Visual representation of amenities on a map interface.
- **Automated Annotations**: Uses machine learning to automate the annotation of geographical data.
- **Deployment Options**: Supports Docker Compose and Kubernetes for deployment.

## 📂 Repository Structure

The file structure is as follows:

```sh
📁
├──📁 Backend # Private and Public Services
│   ├──📁 cmd
│   ├──📁 internal
│   └──📁 sql
├──📁 Distribution
│   ├──📁 compose # Docker Compose Deployment
│   └──📁 kubernetes # Kubernetes Deployment
├──📁 Documentation
│   ├──📁 gantt-chart
│   ├──📁 interim-report
│   ├──📁 mom # Daily Meetings
│   ├──📁 presentation # Weekly Presentation
│   ├──📁 project plan
│   ├──📁 retrospectives # Weekly Retrospectives
│   ├──📁 survey
│   └──📁 ux-documents # General UI/UX Doc Repository
├──📁 Frontend
│   ├──📁 public
│   └──📁 src
└──📁 Python
    ├──📁 charts # Visual Representations
    ├──📁 notebooks # ML Testing and Notes
    ├──📁 script # Production Scripts
    └──📁 yolo # Custom YoLo model
```

## 🚀 Deploying the application

_Magpie_ is provided as a [docker-compose](./Distribution/compose/) and [kubernetes](./Distribution/kubernetes/) deployment.

- The **docker-compose** is intended for local development, but can be used for production if desired (for example on a vps). It's built to be simple, and easy to setup.

- The **kubernetes** deployment is intended for production use. However an existing `load-balancer`, `ingress-controller` and `cert-manager` are required on the cluster.

## ⚙️ Using the development containers

This repository is set up to use dev containers for development. To use them, you will need:

- to have [docker](https://www.docker.com/) installed on your machine
- to have the [dev container extension](https://marketplace.visualstudio.com/items?itemName=ms-vscode-remote.remote-containers) installed in [VS Code](https://code.visualstudio.com/)

Once installed, open the command pallet with `Ctrl+Shift+P` and type `Dev Containers: Open Folder in Container`. Then select the desired folder.

## 📚 Additional Information

Each of the top level directories contains a `README.md` file with more information about the contents of that directory. Please refer to these files for more information about each part of the project.

## 🤝 Contribution

This repo is intended for our Masters Group Project as part of the [TUDublin](https://www.tudublin.ie/) MSc [ASD](https://www.tudublin.ie/study/postgraduate/courses/computing-advanced-software-development-tu059/) / [DS](https://www.tudublin.ie/study/postgraduate/courses/computing-data-science/) course. We are not accepting contributions at this time.

## 📝 License

This project is licensed under the GNU General Public License v3.0 - see the [LICENSE](./licence) file for details.

## 🎖️ Acknowledgments

| Name             | GitHub Username                                        | Email                                                 |
| ---------------- | ------------------------------------------------------ | ----------------------------------------------------- |
| Anais Blenet     | [@anaisbl](https://github.com/anaisbl)                 | TODO                                                  |
| Andreas Kraus    | [@ankraus](https://github.com/ankraus)                 | TODO                                                  |
| Jessica Fornetti | [@JessicaFornetti](https://github.com/JessicaFornetti) | TODO                                                  |
| Kaustubh Trivedi | [@KaustubhTrivedi](https://github.com/KaustubhTrivedi) | TODO                                                  |
| Saul Burgess     | [@1Solon](https://github.com/1Solon)                   | [burgesssaul@gmail.com](mailto:Burgesssaul@gmail.com) |
| Yuanshuo Du      | [@YuanshuoDu](https://github.com/YuanshuoDu)           | TODO                                                  |

_Thanks to our lecturers and (especially) our supervisor for their guidance and support throughout the project!_
