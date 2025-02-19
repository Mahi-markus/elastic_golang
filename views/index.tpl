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
            <div class="relative">
                <input type="text" name="query" id="search-input" placeholder="Search for a product..." class="p-2 w-full border rounded-md" required>
                <!-- Autocomplete suggestions box -->
                <div id="autocomplete-suggestions" class="absolute w-full bg-white border rounded-md mt-1 hidden z-10"></div>
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
            
            // Make AJAX request for the selected query to fetch detailed information
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

        // Handle autocomplete suggestions as the user types
        document.getElementById('search-input').addEventListener('input', function() {
            let query = this.value;
            let suggestionBox = document.getElementById('autocomplete-suggestions');

            if (query.length < 2) {
                suggestionBox.innerHTML = '';
                suggestionBox.classList.add('hidden');
                return;
            }

            // Fetch autocomplete suggestions from the backend (you should implement this API on the backend)
            fetch(`/autocomplete?query=${query}`, {
                method: 'GET',
                headers: {
                    'Content-Type': 'application/json',
                }
            })
            .then(response => response.json())
            .then(data => {
                suggestionBox.innerHTML = ''; // Clear the suggestion box
                if (data.length > 0) {
                    data.forEach(product => {
                        let suggestionItem = document.createElement('div');
                        suggestionItem.textContent = product;
                        suggestionItem.classList.add('p-2', 'cursor-pointer');
                        suggestionItem.addEventListener('click', function() {
                            document.getElementById('search-input').value = product;
                            suggestionBox.classList.add('hidden');
                            // Trigger search
                            document.getElementById('search-form').submit();
                        });
                        suggestionBox.appendChild(suggestionItem);
                    });
                    suggestionBox.classList.remove('hidden'); // Show suggestion box
                } else {
                    suggestionBox.classList.add('hidden'); // Hide if no suggestions
                }
            })
            .catch(error => {
                console.error('Error fetching autocomplete suggestions:', error);
                suggestionBox.classList.add('hidden'); // Hide suggestion box on error
            });
        });
    </script>
</body>
</html>
