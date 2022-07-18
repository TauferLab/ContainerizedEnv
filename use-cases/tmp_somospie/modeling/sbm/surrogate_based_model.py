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

import scipy.stats as stats
import sys
import argparse
import numpy as np
import numpy.matlib
import random
import pandas as pd
from itertools import combinations_with_replacement as cwr 
from itertools import chain
import matplotlib.pyplot as plt
from matplotlib.colors import ListedColormap
cmap_brg = plt.cm.get_cmap('brg', 256)
clist_rg = cmap_brg(np.linspace(0.5, 1, 128))
cmap_rg = ListedColormap(clist_rg)
plt.rcParams["figure.figsize"] = (20,10)

############################################################################################################################################
def get_random_sample(original_values,N_samples):
    N = original_values.shape[0]
    indices = list(range(N))
    Rand_smpl = sorted(random.sample(indices, N_samples))
    sampled_values =  original_values.iloc[Rand_smpl]
    return sampled_values

############################################################################################################################################


### This function return the binomial coefficient n Choose r.
def nCr(n, r):
	if n < r:
		return 0
	retVal = 1
	for i in range(r):
		retVal = (retVal*(n-i))/(i+1)
	
	return int(retVal)

############################################################################################################################################

### Generate the random partition of points into the K-folds

def K_partition(N, NUM_FOLDS):
    random.seed(None)
    indices = list(range(N))
    Folds = []
    #print('random seed in K_partition', random.getstate()[1][0])

    
    ### Generate the random partition of points into the K-folds
    for i in range(NUM_FOLDS-1):
        Rand_smpl = sorted(random.sample(indices, int(N/NUM_FOLDS)))
        #print('Random samples inside the code', Rand_smpl)
        Folds.append(Rand_smpl)

        for j in Rand_smpl:
            indices.remove(j)

    Folds.append(indices)
    return Folds

############################################################################################################################################

### Generate the random partition of points into the K-folds
def K_folds(original_values,NUM_FOLDS,degrees):
        ### Seed the random number generator, None causes the seed to be the system time.
    random.seed(None)
    #print('random seed', random.getstate()[1][0])

    N = original_values.shape[0]			### Number of data points in the file

    indices = list(range(N))

    Folds = []

    ### Generate the random partition of points into the K-folds
    Folds = K_partition(N, NUM_FOLDS)

    #SAVE_SSE = []
    Kfold_parameters = np.matlib.zeros((int(len(degrees)), 2))
    # Extracting the training and the testing values
    indice_Kfold = 0
    for degree in degrees:
        SSE = 0 
        for i in range(NUM_FOLDS):
          
            training = np.delete(original_values, Folds[i], axis=0) # Delete all the rows in Folds[i]
            # dependen variable
            dependent_variable_values = training[:,-1]
            
            # Independent Variables
            independent_variable_points = training[:,:-1]
            
            
            
            Parameters = int(independent_variable_points.shape[1])


            nMonomials = nCr(degree+Parameters, Parameters)	### This counts the number of monomials that will appear in the polynomial of degree DEGREE in two variables.
            N = independent_variable_points.shape[0]					### The number of data points
            ### These are the matrices X, Z, and $\beta$ as defined in the paper:
            ### "Performance Tuning of MapReduce Jobs Using Surrogate-Based Modeling" by Travis Johnston, Mohammad Alsulmi, Pietro Cicotti, and Michela Taufer.
            ### The description is found in Section 2.2, equations (2).
            X = np.matlib.ones((int(N), int(nMonomials)))
            Beta = np.matlib.zeros((int(nMonomials), 1))
            Z = np.matlib.zeros((int(N), 1))


            ### Construct the values in matrix X and Z.
            X = [ [np.product(x) for x in cwr(chain([1.0], iv), degree)] for iv in independent_variable_points ]
            X =  np.array(X)
            Z = np.array(dependent_variable_values)
            Zbar = np.mean(Z)

            ### Perform the multiplication that gives us the Beta vector (of coefficients)
            ### We want to compute Beta = (X^T X)^-1 X^T Z.
            ### We compute it in two steps.
            X_transX = np.linalg.inv(np.matmul(X.transpose(),X))	### X_transX = (X^T X)^-1 and is computed separately since we will recycle this portion over and over again later.
            Beta = np.matmul(X_transX,np.matmul(X.transpose(),Z))

            testing  = original_values[Folds[i]]
            #testing_values = np.array(testing[[testing.columns[-1]]])
            testing_values = testing[:,-1]
            # Independent Variables
            #Evaluation_values = testing.drop(testing.columns[-1], axis=1).values
            Evaluation_values = testing[:,:-1]
            Evaluation = [evaluate_polynomial(Beta, degree, point) for point in Evaluation_values]
            Evaluation = np.array([Evaluation])
            Evaluation = Evaluation.transpose()

            # Difference of the two vectors
            diffs = testing_values - Evaluation
            SSE = SSE + np.linalg.norm(diffs,2)**2
        #SSE_DF = {'Degree':degree, 'SSE':SSE/NUM_FOLDS}
        #Kfold_parameters = Kfold_parameters.append(SSE_DF, ignore_index=True)
        Kfold_parameters[indice_Kfold,0] = degree
        Kfold_parameters[indice_Kfold,1] = SSE
        indice_Kfold = indice_Kfold+1
    return Kfold_parameters

############################################################################################################################################

def find_best_degree(original_values, epochs, k_folds):
    # Finding the best degree
    degrees = list(range(1,10))
    Kfold_parameter = np.zeros((1,2))
    for i in range(epochs):
        print(f'Training model')
        print(f'I am in the epoch {i}')
        temp = K_folds(original_values,k_folds,degrees)
        Kfold_parameter = np.append(Kfold_parameter,temp, axis=0)
    Kfold_parameter = Kfold_parameter[1:] 


    # # Saving the information of the data with the best degree
    indice = np.argmin(Kfold_parameter[:,1], axis=0)
    indice = indice[0,0]
    degree = int(Kfold_parameter[indice,0])

    np.argmin(Kfold_parameter[:,1], axis=0)
    return degree
    
############################################################################################################################################

def construct_model(original_values, epochs, k_folds):
    # Find the best degree
    degree = find_best_degree(original_values, epochs, k_folds)

    # Dependent variable
    dependent_variable_values = original_values[:,-1]

    # Independent Variables
    independent_variable_points = original_values[:,:-1]
     
	### independent_variable_points: independent values of training. Numpy array
	### dependent_variable_values  : expected values in training file. Numpy array
	### degree                     : polynomial degree of the Modeling. Int
    Parameters = int(independent_variable_points.shape[1])
    nMonomials = nCr(int(degree)+int(Parameters), int(Parameters))	### This counts the number of monomials that will appear in the polynomial of degree DEGREE in two variables.
    N = independent_variable_points.shape[0]					### The number of data points
    df = N - nMonomials - 1			### Number of degrees of freedom, a parameter to pass to stats.t

    ### These are the matrices X, Z, and $\beta$ as defined in the paper:
    ### "Performance Tuning of MapReduce Jobs Using Surrogate-Based Modeling" by Travis Johnston, Mohammad Alsulmi, Pietro Cicotti, and Michela Taufer.
    ### The description is found in Section 2.2, equations (2).
    X = np.matlib.ones((int(N), int(nMonomials)))
    Beta = np.matlib.zeros((int(nMonomials), 1))
    Z = np.matlib.zeros((int(N), 1))


    ### Construct the values in matrix X and Z.
    X = [ [np.product(x) for x in cwr(chain([1.0], iv), degree)] for iv in independent_variable_points ]
    X =  np.array(X)
    Z = np.array(dependent_variable_values)
    Zbar = np.mean(Z)

    ### Perform the multiplication that gives us the Beta vector (of coefficients)
    ### We want to compute Beta = (X^T X)^-1 X^T Z.
    ### We compute it in two steps.
    X_transX = np.linalg.inv(np.matmul(X.transpose(),X))	### X_transX = (X^T X)^-1 and is computed separately since we will recycle this portion over and over again later.
    Beta = np.matmul(X_transX,np.matmul(X.transpose(),Z))
    return Beta, degree

############################################################################################################################################

# This function expects:
# * a list of coefficients for the polynomial in order: 
# * the degree of the polynomial (integer).
# * a point (list of floats) of where to evaluate the polynomial.
# This function returns the value of the polynomial evaluated at the point provided.
def evaluate_polynomial(coefficients, degree, point):
    if degree == 0:
        return coefficients[0]
    
    monomials = [ np.product(x) for x in cwr(chain([1.0], point), degree) ]
    return sum( [ a[0]*a[1] for a in zip(coefficients, monomials) ] )


############################################################################################################################################

def predictions_to_csv(predictions, output_data):
    predictions = predictions[['x','y','sm']]
    # Saving the predicted values to a csv document
    predictions.to_csv(output_data, index=False)
    
############################################################################################################################################

def evaluate_model(Evaluation_DF, Beta, degree, output_data):
    
    Evaluation_values = Evaluation_DF.values
    Evaluation = [evaluate_polynomial(Beta, degree, point) for point in Evaluation_values]
    Evaluation = pd.DataFrame(Evaluation,columns=['sm'])
    # Writing in the Evaluation file
    Evaluation = Evaluation_DF.join(Evaluation)
   
    #Print predictions to csv
    predictions_to_csv(Evaluation, output_data)
    
############################################################################################################################################

def reducing_neigbors(original_values,degree,Evaluation):
    # original values to create the model
    # degree of the polynomial
    # Exhaustive predicted values to compare with the model created with fewer points
    
    # Compute the minimum number of points to create the model;
    parameters   = original_values.shape[1]-1 # Number of parameters
    N            = original_values.shape[0]   # Number of points in the data
    nPoints = nCr(degree+parameters, parameters) # minimum nunber of points to build the model
    
    Extra_points = min(int((N- nPoints)/2),150)
    
    # Randomly extract "nPoints+Extra_points" from the original data
    sampled_values = get_random_sample(original_values,nPoints+Extra_points)
    
    exhaustive_variable = Evaluation[[Evaluation.columns[-1]]]
    exhaustive_variable_vector = exhaustive_variable.values
    # Independent Variables
    independent_variable = Evaluation.drop(Evaluation.columns[-1], axis=1)
    
    K_parameters = pd.DataFrame()
    SAVE_SSE = []
    
    for i in range(1):
        extra_point = N - nPoints
        sampled_points = sampled_values.iloc[:(nPoints+extra_point)]                             # Extract the points
        Beta = construct_model(sampled_points, degree)
        Beta = Beta.transpose()
        Beta = list(Beta)  
        K_parameters = K_parameters.append(Beta, ignore_index=True)
        
        # Build the model  
        Evaluation_temp = evaluate_model(independent_variable, Beta, degree)                     # Evaluate the model
        # Build the model
        Temp_variable = Evaluation_temp[[Evaluation_temp.columns[-1]]]
        Temp_variable_vector = Temp_variable.values
        diffs =  exhaustive_variable_vector - Temp_variable_vector
        SSE = np.linalg.norm(diffs,2)**2
        SAVE_SSE.append(SSE)

        
        
    SAVE_SSE = pd.DataFrame(SAVE_SSE)
    SAVE_SSE.columns = ['SSE']   
    Beta_columns = [f'B{i}' for i in range(int(nPoints))] 
    K_parameters.columns = Beta_columns
    K_parameters = SAVE_SSE.join(K_parameters)
    
    return K_parameters

##########################################################################################################################################

# General heatmap function. 
def heatmap(df, horizontal=0, vertical=1, value=2, vmin=None, vmax=None, size=1, title="Heatmap", out="", cmap=None):
    fig, ax = plt.subplots()
    df.plot.scatter(x=horizontal, y=vertical, s=size, c=value, cmap=cmap, title=title, vmin=vmin, vmax=vmax, ax=ax)
    if out:
        print(f"Saving image to {out}")
        plt.savefig(out)
        plt.xlabel('Longitude', fontsize=14)
        plt.ylabel('Latitude', fontsize=14)
        plt.show()
    else:
        plt.xlabel('Longitude', fontsize=14)
        plt.ylabel('Latitude', fontsize=14)
        plt.show()
    plt.close()
    
############################################################################################################################################

def sbm_matrix(original_values, degree):
    # Dependent variable
    
    dependent_variable_values = original_values[[original_values.columns[-1]]].values

    # Independent Variables
    independent_variable_points = original_values.drop(original_values.columns[-1], axis=1).values
    
	### independent_variable_points: independent values of training. Numpy array
	### dependent_variable_values  : expected values in training file. Numpy array
	### degree                     : polynomial degree of the Modeling. Int
    Parameters = int(independent_variable_points.shape[1])
    nMonomials = nCr(int(degree)+int(Parameters), int(Parameters))	### This counts the number of monomials that will appear in the polynomial of degree DEGREE in two variables.
    N = independent_variable_points.shape[0]					### The number of data points
    df = N - nMonomials - 1			### Number of degrees of freedom, a parameter to pass to stats.t

    ### These are the matrices X, Z, and $\beta$ as defined in the paper:
    ### "Performance Tuning of MapReduce Jobs Using Surrogate-Based Modeling" by Travis Johnston, Mohammad Alsulmi, Pietro Cicotti, and Michela Taufer.
    ### The description is found in Section 2.2, equations (2).
    X = np.matlib.ones((int(N), int(nMonomials)))
    Beta = np.matlib.zeros((int(nMonomials), 1))
    Z = np.matlib.zeros((int(N), 1))


    ### Construct the values in matrix X and Z.
    X = [ [np.product(x) for x in cwr(chain([1.0], iv), degree)] for iv in independent_variable_points ]
    X =  np.array(X)

    return X

############################################################################################################################################

def indices_of_kNN( scaled_data, specific_point, k):

    #scaled_data = data_points[[data_points.columns[-1]]].values
    
    distances = [ sum( (x-specific_point)**2 ) for x in scaled_data ]
    indices = np.argsort( distances, kind='mergesort' )[:k+1]

    if distances[indices[0]] < .001:
        return indices[1:]
    else:
        return indices[:k]
