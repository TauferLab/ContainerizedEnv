#!/usr/bin/env python3
# coding: utf-8

# Copyright (c)
# 2015 by The University of Delaware
# 2020 by The University of Tennessee Knoxville
# Contributors: Travis Johnston, Leobardo Valera, and Ricardo Llamas
# Affiliation: Global Computing, Laboratory, Michela Taufer PI
# Url: http://gcl.cis.udel.edu/, https://github.com/TauferLab
# 
# All rights reserved.
# 
# Redistribution and use in source and binary forms, with or without modification,
# are permitted provided that the following conditions are met:
#  
# 1. Redistributions of source code must retain the above copyright notice, 
# this list of conditions and the following disclaimer.
# 
# 2. Redistributions in binary form must reproduce the above copyright notice,
# this list of conditions and the following disclaimer in the documentation
# and/or other materials provided with the distribution.
# 
# 3. If this code is used to create a published work, one or both of the
# following papers must be cited.
# 
# Travis Johnston, Mohammad Alsulmi, Pietro Cicotti, Michela Taufer, "Performance
# Tuning of MapReduce Jobs Using Surrogate-Based Modeling," International
# Conference on Computational Science (ICCS) 2015.
# 
# M. Matheny, S. Herbein, N. Podhorszki, S. Klasky, M. Taufer, "Using Surrogate-
# based Modeling to Predict Optimal I/O Parameters of Applications at the Extreme
# Scale," In Proceedings of the 20th IEEE International Conference on Parallel and
# Distributed Systems (ICPADS), December 2014.
# 
# 4.  Permission of the PI must be obtained before this software is used
# for commercial purposes.  (Contact: taufer@acm.org)
# 
# THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
# ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED 
# WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED.
# IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT,
# INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING,
# BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
# DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF
# LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE
# OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED
# OF THE POSSIBILITY OF SUCH DAMAGE.
#
import sys
import numpy as np
import pandas as pd
import argparse
sys.path.insert(0, "code/modeling/sbm")
import surrogate_based_model as sbm


#Translate from namespaces to Python variables 
def from_args_to_vars (args):
    print("Reading training data from", args.trainingdata)
    training_data = pd.read_csv(args.trainingdata)
    sm = training_data[['z']]
    training_data = training_data.drop(['z'], axis=1)
    # Move the column of Moisture to the end
    training_data = pd.concat([training_data, sm], axis=1)
    # Convert the Data Frame to a matrix
    training_data = training_data.values
    
    epochs = int(args.epochs)
    k_folds = int(args.kfolds)
    evaluation_data = pd.read_csv(args.evaluationdata) 
    output_data = args.outputdata
    return training_data, evaluation_data, output_data, epochs, k_folds

if __name__ == "__main__":	
    #Input arguments to execute the sbm 
    parser = argparse.ArgumentParser(description='Arguments and data files for executing Nearest Neighbors Regression.')
    parser.add_argument('-t', "--trainingdata", help='Training data')
    parser.add_argument('-e', "--evaluationdata", help='Evaluation data')
    parser.add_argument('-o', "--outputdata", help='Predictions')
    parser.add_argument('-epochs', "--epochs", help='The number of times the algorithm sees the entire data set', default=2)
    parser.add_argument('-kfolds', "--kfolds", help='NUmber of folds for hyperparameters search', default=10)
    args = parser.parse_args()

    training_data, evaluation_data, output_data, epochs, k_folds = from_args_to_vars(args)
    # Finding the coefficients of the surface with the minimum SSE
    beta, degree = sbm.construct_model(training_data, epochs, k_folds)
    # Evaluating the model (predicted values)
    sbm.evaluate_model(evaluation_data, beta, degree, output_data)
