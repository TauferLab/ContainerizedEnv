# Traceability and Reproducibility through Individual Containerization (TRIC) 

This repository contains a computational containerized environment that automatically creates a workflow execution's record trail and invisibly attaches it to the workflow's output, enabling data traceability and results explainability.

**Motivation:** To trust findings in computational science, scientists need workflows that trace the data provenance and support results explainability. As workflows become more complex, tracing data provenance and explaining results become harder to achieve. Our solution transforms existing container technology, includes tools for automatically annotating provenance metadata, and allows effective movement of data and metadata across the workflow execution.   

**Environment overview:** Our environment includes two stages: the execution of the workflow while automatically collecting metadata, and the analysis of the metadata. For the first stage, our environment decouples data and applications of traditionally tightly coupled workflows and encapsulates them into individual fine-grained containers. We augment both data and application containers to expose provenance metadata and to move data across the containerized workflow effectively. For the second stage, we provide an interface for visualizing and studying the metadata that scientists can use to understand the data lineage and the computational methods.

**Use case:** We include the demonstration of the capabilities of our environment with the study of an earth science workflow. This workflow predicts soil moisture values from the 27 km resolution satellite data down to the fine-grain 10 m resolution necessary for practical use in policymaking and precision agriculture using a suite of machine learning modeling techniques. By running the workflow in our environment, the end-user can identify the causes of different accuracy for predicted soil moisture values in different resolutions of the input data, and link different results to different machine learning methods used during the soil moisture downscaling, all without requiring the scientist to know aspects of workflow design and implementation. 

This document is organized in the following order:
* [Installation](#installation)
* [Workflow setup guide](#workflow-setup-guide)
* [Metadata interface guide](#metadata-interface-guide)
* Additional project details
  1. [Acknowledgments](#acknowledgments)
  2. [The project team](#project-team)
* [Copyright and license information](#copyright-and-license)

## Installation 
There are two main components to install to run our environment: 1) Apptainer and the plugins for the execution of the workflow and automatic collection of the metadata; and 2) a Jupyter Notebook for the interface that enables the visualization and study of the metadata to understand the data lineage and the computational methods.

### Apptainer and plugins for containerized workflow creation and execution 
1. Install Apptainer by following these [instructions.](https://github.com/apptainer/apptainer/blob/main/INSTALL.md) A version greater than 3.5 is required to enable the plugins and zero-copy data transfer between containers.   

2. Clone our repository
    ```
    git clone --recurse-submodules https://github.com/TauferLab/Src_ContainerizedEnv
    ```
3. Install the metadata plugin  
    ```
    apptainer plugin compile plugin/.  
    sudo plugin install plugin/plugin.sif 
    ```
### Interface for metadata analysis 
The interface is a Jupyter notebook that has the next required dependencies:  
* python=3.9.12
* matplotlib=3.5.1
* pip=22.1.2
* jupyter
* ipyfilechooser
* networkx=2.8.4
* rglob=1.7
* sif=0.0.11
* networkx_viewer=0.3.0 
  
We use anaconda to install the software stack. If you do not have Anaconda installed, you can follow the instructions [here](https://docs.anaconda.com/anaconda/install/index.html) to install it.
Once you have anaconda, you can create the environment. To this end, **make sure to change the prefix in `install/env_conda.yml` to the location of Anaconda in your local machine (e.g., `/opt/anaconda3/`, `/home/opt/anaconda3/`)**. You can use `whereis conda` to check the path. 

Run the next commands on your local machine:  
```
conda env create -f install/env_conda.yml
conda activate tric
```
Once you have your environment installed, you can run `jupyter notebook` and run the `interface/metadata-interface.ipynb`

## Workflow setup guide 

1. Pick a sample workflow  
    * For knn: `cd sample_workflows/knn`  
    * For sbm: `cd sample_workflows/sbm`  
    * For rf:  `cd sample_workflows/rf`  

2. Initilize the workflow containers  
    * For knn: `apptainer workflow --create knn_workflow.json`  
    * For sbm: `apptainer workflow --create sbm_workflow.json`  
    * For rf:  `apptainer workflow --create rf_workflow.json`  
    * Or create your own workflow: `apptainer workflow --create`

3. Run the workflow  
    * For knn: `apptainer workflow --run knn_workflow.json`  
    * For sbm: `apptainer workflow --run sbm_workflow.json`  
    * For rf:  `apptainer workflow --run rf_workflow.json`  
    * Or run your own workflow: `apptainer workflow --run your_workflow_name.json`

4. Explore the metadata using the metadata interface  

## Metadata interface guide  

1. Navigate to your desired metadata directory
    * For knn: `cd sample_workflows/knn/metadata`  
    * For sbm: `cd sample_workflows/sbm/metadata`  
    * For rf:  `cd sample_workflows/rf/metadata`  

2. Use the metadata interface  
    * Run `jupyter notebook` and select the `interface/metadata-interface.ipynb` notebook
    * Execute the first cell to import all libraries and packages
    * Execute the second cell to select one file or a directory with metadata files. For example, select `Src_ContainerizedEnv/sample_workflows/knn/metadata/` 
    * Execute the third cell to build a graph representation of the workflow
    * Execute the fourth cell to open a visualizer where you can interact with each of the components and observe the metadata 

 ## Additional project details
 
### Acknowledgments

We aknowledge the support of Sandia National Laboratories; the National Science Foundation through the awards 1841758, 2103845, 2138811, and 1941443; and IBM through a Shared University Research Award. The authors acknowledge the Singularity team, specially Cedric Clerget and Ian Kaneshiro, for the support.
This work was partially developed and tested using Jeststream [XSEDE computing resources](https://portal.xsede.org/group/xup/resource-monitor). 

### Project team

Developers:
* Paula Olaya
* Dominic Kennedy 

Project Advisors:
* Dr. Jay Lofstead
* Dr. Michela Taufer 


## Copyright and license

Copyright (c) 2022, Global Computing Lab
