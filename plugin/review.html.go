package main

const reviewHTML = `
<!doctype html>
<html lang="en">
    <head>
        <meta charset="utf-8">
        <meta name="viewport" content="width=device-width, initial-scale=1">

        <!-- Bootstrap CSS -->
        <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.0.2/dist/css/bootstrap.min.css" rel="stylesheet" integrity="sha384-EVSTQN3/azprG1Anm3QDgpJLIm9Nao0Yz1ztcQTwFspd3yD65VohhpuuCOmLASjC" crossorigin="anonymous">

        <title>Workflow review</title>
    </head>
    <body>
        <nav class="navbar navbar-expand-md navbar-dark bg-dark">
            <div class="container-fluid">
                <a class="navbar-brand" href="https://globalcomputing.group/research.html">Workflow creation wizard (TRIC): Review</a>
                <div class="navbar-nav ms-auto mb-2 mb-md-0">
                    <!--<div class="nav-item">-->
                    <!--<a class="nav-link" href="https://www.ospray.org/">About</a>-->
                    <!--</div>-->
                </div>
            </div>
        </nav>

        <div class="container-lg pt-3 pb-3">

            <div class="card mb-3">
                <div class="card-body">
                    <h5 class="card-title">Application Container</h5>
                    <p class="card-text">
                    <p> Name: {{.ApplicationContainer.Name}} </p>
                    <p> Path: {{.ApplicationContainer.InPath}} </p>
                    </p>
                </div>
            </div>

            <div class="card mb-3">
                <div class="card-body">
                    <h5 class="card-title">Input Containers</h5>
                    {{range .InputContainer}}
                    <div class="card mb-2">
                        <div class="card-body">
                            <p class="card-text">
                            <p> Name: {{.Name}} </p>
                            <p> Path: {{.InPath}} </p>
                            <p> Size: {{.Size}} </p>
                            </p>
                        </div>
                    </div>
                    {{end}}
                </div>
            </div>

            <div class="card mb-3">
                <div class="card-body">
                    <h5 class="card-title">Ouptput Container</h5>
                    <p class="card-text">
                    <p> Name: {{.OutputContainer.Name}} </p>
                    <p> Size: {{.OutputContainer.Size}} </p>
                    </p>
                </div>
            </div>

            <form action="/quit">
                <button class="btn btn-dark">Quit</button>
            </form>
        </div>


        <!-- Option 1: Bootstrap Bundle with Popper -->
        <!--<script src="https://cdn.jsdelivr.net/npm/bootstrap@5.0.2/dist/js/bootstrap.bundle.min.js" integrity="sha384-MrcW6ZMFYlzcLA8Nl+NtUVF0sA7MsXsP1UyJoMp4YLEuNSfAP+JcXn/tWtIaxVXM" crossorigin="anonymous"></script>-->

        <!-- Option 2: Separate Popper and Bootstrap JS -->
        <!--
            <script src="https://cdn.jsdelivr.net/npm/@popperjs/core@2.9.2/dist/umd/popper.min.js" integrity="sha384-IQsoLXl5PILFhosVNubq5LC7Qb9DXgDA9i+tQ8Zj3iwWAwPtgFTxbJ8NT4GN1R8p" crossorigin="anonymous"></script>
            <script src="https://cdn.jsdelivr.net/npm/bootstrap@5.0.2/dist/js/bootstrap.min.js" integrity="sha384-cVKIPhGWiC2Al4u+LWgxfKTRIcfu0JTxR+EQDz/bgldoEyl4H0zUF0QKbrJ0EcQF" crossorigin="anonymous"></script>
            -->
    </body>
</html>
`
