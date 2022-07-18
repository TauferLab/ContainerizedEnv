package main

const indexHTML = `
<!doctype html>
<html lang="en">
    <head>
        <meta charset="utf-8">
        <meta name="viewport" content="width=device-width, initial-scale=1">

        <!-- Bootstrap CSS -->
        <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.0.2/dist/css/bootstrap.min.css" rel="stylesheet" integrity="sha384-EVSTQN3/azprG1Anm3QDgpJLIm9Nao0Yz1ztcQTwFspd3yD65VohhpuuCOmLASjC" crossorigin="anonymous">

        <title>TRIC</title>
        <script>
            function appendInputContainer() {
                            const node = document.getElementById("inputContList");
                            let inputContTemp = document.getElementById("inputContTemp");
                            node.appendChild(inputContTemp.content.cloneNode(true));
                        }
        </script>
    </head>
    <body>
        <nav class="navbar navbar-expand-md navbar-dark bg-dark">
            <div class="container-fluid">
                <a class="navbar-brand" href="https://globalcomputing.group/research.html">Workflow creation wizard (TRIC)</a>
                <div class="navbar-nav ms-auto mb-2 mb-md-0">
                </div>
            </div>
        </nav>

        <form method="POST" action="/post">     
            <div class="container-lg pt-3">

                <div class="card mb-3">
                    <div class="card-body">
                        <h5 class="card-title">Application Container</h5>
                        <p class="card-text">
                        <div class="form-floating mb-2">
                            <input type="text" class="form-control" id="applicationContainer.name" name="applicationContainer.name" placeholder="applicationContainer.name">
                            <label for="applicationContainer.name">Application container name</label>
                        </div>
                        <div class="form-floating mb-2">
                            <input type="text" class="form-control" id="applicationContainer.inPath" name="applicationContainer.inPath" placeholder="applicationContainer.inPath">
                            <label for="applicationContainer.inPath">Application container definition file path</label>
                        </div>
                        </p>
                    </div>
                </div>

                <div class="card mb-3">
                    <div class="card-body">
                        <h5 class="card-title">Input Containers</h5>
                        <p class="card-text">
                        <button type="button" class="btn btn-secondary" onclick="appendInputContainer()">Add input container</button>
                        <div id="inputContList"></div>
                        </p>
                    </div>
                </div>

                <template id="inputContTemp">
                    <div class="card mb-3">
                        <div class="card-body">
                            <h5 class="card-title">Input Container</h5>
                            <p class="card-text">
                            <div class="form-floating mb-2">
                                <input type="text" class="form-control" id="inputContainer.name" name="inputContainer.name" placeholder="inputContainer.name">
                                <label for="inputContainer.name">Input container name</label>
                            </div>
                            <div class="form-floating mb-2">
                                <input type="text" class="form-control" id="inputContainer.inPath" name="inputContainer.inPath" placeholder="inputContainer.inPath">
                                <label for="inputContainer.inPath">Input data path (optional)</label>
                            </div>
                            <div class="input-group mb-2">
                                <label class="input-group-text" for="inputContainer.size">Input container size</label>
                                <input type="number" class="form-control" id="inputContainer.size" name="inputContainer.size">
                                <select class="form-select" name="inputContainer.sizeUnit">
                                    <option value="0">Bytes</option>
                                    <option value="1">Kilobytes</option>
                                    <option value="2">Megabytes</option>
                                    <option value="3">Gigabytes</option>
                                    <option value="4">Terabytes</option>
                                    <option value="5">Petabytes</option>
                                </select>
                            </div>
                            </p>
                        </div>
                    </div>
                </template>

                <div class="card mb-3">
                    <div class="card-body">
                        <h5 class="card-title">Output Container</h5>
                        <p class="card-text">
                        <div class="form-floating mb-2">
                            <input type="text" class="form-control" id="outputContainer.name" name="outputContainer.name" placeholder="outputContainer.name">
                            <label for="outputContainer.name">Output container name</label>
                        </div>
                        <div class="input-group mb-2">
                            <label class="input-group-text" for="outputContainer.size">Output container size</label>
                            <input type="number" class="form-control" id="outputContainer.size" name="outputContainer.size">
                            <select class="form-select" name="outputContainer.sizeUnit">
                                <option value="0">Bytes</option>
                                <option value="1">Kilobytes</option>
                                <option value="2">Megabytes</option>
                                <option value="3">Gigabytes</option>
                                <option value="4">Terabytes</option>
                                <option value="5">Petabytes</option>
                            </select>
                        </div>
                        </p>
                    </div>
                </div>

                <div class="card mb-3">
                    <div class="card-body">
                        <h5 class="card-title">Workflow name</h5>
                        <p class="card-text">
                        <div class="form-floating mb-2">
                            <input type="text" class="form-control" id="workflowName" name="workflowName" placeholder="workflowName">
                            <label for="workflowName">Workflow name</label>
                        </div>
                        <button class="btn btn-primary">Create workflow</button>
                        <a href="quit" class="btn btn-dark">Quit</a>
                        </p>
                    </div>
                </div>
            </div>
        </form>

        <!--<script src="https://cdn.jsdelivr.net/npm/bootstrap@5.0.2/dist/js/bootstrap.bundle.min.js" integrity="sha384-MrcW6ZMFYlzcLA8Nl+NtUVF0sA7MsXsP1UyJoMp4YLEuNSfAP+JcXn/tWtIaxVXM" crossorigin="anonymous"></script>-->

        <!--
            <script src="https://cdn.jsdelivr.net/npm/@popperjs/core@2.9.2/dist/umd/popper.min.js" integrity="sha384-IQsoLXl5PILFhosVNubq5LC7Qb9DXgDA9i+tQ8Zj3iwWAwPtgFTxbJ8NT4GN1R8p" crossorigin="anonymous"></script>
            <script src="https://cdn.jsdelivr.net/npm/bootstrap@5.0.2/dist/js/bootstrap.min.js" integrity="sha384-cVKIPhGWiC2Al4u+LWgxfKTRIcfu0JTxR+EQDz/bgldoEyl4H0zUF0QKbrJ0EcQF" crossorigin="anonymous"></script>
            -->
    </body>
</html>
`
