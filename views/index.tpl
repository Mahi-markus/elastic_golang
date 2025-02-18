<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Product Search</title>
    <link href="https://cdn.jsdelivr.net/npm/tailwindcss@2.2.19/dist/tailwind.min.css" rel="stylesheet">
</head>
<body class="bg-gray-100 p-4">
    <div class="container mx-auto">
        <h1 class="text-2xl font-semibold mb-6">Search Products</h1>
        
        <form id="search-form" action="/search" method="get" class="mb-4">
            <div class="flex items-center">
                <input type="text" name="query" id="search-input" placeholder="Search for a product..." class="p-2 w-full border rounded-md" required>
                <button type="submit" class="ml-2 p-2 bg-blue-500 text-white rounded-md">Search</button>
            </div>
        </form>

        <!-- Display search results -->
        <div id="results" class="mt-4">
            <!-- Results will appear here -->
        </div>
    </div>

    <script>
        // Handle form submission using AJAX to get search results without reloading the page
        document.getElementById('search-form').addEventListener('submit', function(e) {
            e.preventDefault();
            let query = document.getElementById('search-input').value;
            
            // Make AJAX request
            fetch(`/search?query=${query}`, {
                method: 'GET',
                headers: {
                    'Content-Type': 'application/json',
                }
            })
            .then(response => response.json())
            .then(data => {
                // Clear previous results
                let resultsDiv = document.getElementById('results');
                resultsDiv.innerHTML = '';

                // Display results
                if (data.length > 0) {
                    let list = '<ul>';
                    data.forEach(product => {
                        list += `<li class="p-2 border-b">${product}</li>`;
                    });
                    list += '</ul>';
                    resultsDiv.innerHTML = list;
                } else {
                    resultsDiv.innerHTML = '<p>No results found.</p>';
                }
            })
            .catch(error => {
                console.error('Error:', error);
            });
        });
    </script>
</body>
</html>
