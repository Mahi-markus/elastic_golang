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
                <div id="autocomplete-suggestions" class="absolute w-full bg-white border rounded-md mt-1 hidden z-10"></div>
            </div>
        </form>

        <div id="results" class="mt-4"></div>

        <div id="product-details" class="mt-6 p-4 border bg-white rounded-md hidden">
            <h2 class="text-xl font-semibold">Product Details</h2>
            <div id="product-name" class="mt-2 text-lg font-bold"></div>
            <div id="product-description" class="mt-2"></div>
            <div id="product-price" class="mt-2 font-bold"></div>
            <div id="product-manufacturer" class="mt-2"></div>
        </div>
    </div>

    <script>
        document.getElementById('search-form').addEventListener('submit', function(e) {
            e.preventDefault();
            let query = document.getElementById('search-input').value;

            fetch(`/search?query=${query}`, {
                method: 'GET',
                headers: {
                    'Content-Type': 'application/json',
                }
            })
            .then(response => response.json())
            .then(data => {
                if (data.error) {
                    document.getElementById('results').innerHTML = `<p>${data.error}</p>`;
                    document.getElementById('product-details').classList.add('hidden');
                    return;
                }

                document.getElementById('product-name').textContent = `Name: ${data.name}`;
                document.getElementById('product-description').textContent = `Description: ${data.description}`;
                document.getElementById('product-price').textContent = `Price: ${data.base_price}`;
                document.getElementById('product-manufacturer').textContent = `Manufacturer: ${data.manufacturer}`;
                document.getElementById('product-details').classList.remove('hidden');
            })
            .catch(error => {
                console.error('Error:', error);
            });
        });

        document.getElementById('search-input').addEventListener('input', function() {
            let query = this.value;
            let suggestionBox = document.getElementById('autocomplete-suggestions');

            if (query.length < 2) {
                suggestionBox.innerHTML = '';
                suggestionBox.classList.add('hidden');
                return;
            }

            fetch(`/autocomplete?query=${query}`, {
                method: 'GET',
                headers: {
                    'Content-Type': 'application/json',
                }
            })
            .then(response => response.json())
            .then(data => {
                suggestionBox.innerHTML = '';
                if (data.length > 0) {
                    data.forEach(product => {
                        let suggestionItem = document.createElement('div');
                        suggestionItem.textContent = product;
                        suggestionItem.classList.add('p-2', 'cursor-pointer', 'hover:bg-gray-200');
                        suggestionItem.addEventListener('click', function() {
                            document.getElementById('search-input').value = product;
                            suggestionBox.classList.add('hidden');
                            showProductDetails(product);
                        });
                        suggestionBox.appendChild(suggestionItem);
                    });
                    suggestionBox.classList.remove('hidden');
                } else {
                    suggestionBox.classList.add('hidden');
                }
            })
            .catch(error => {
                console.error('Error fetching autocomplete suggestions:', error);
                suggestionBox.classList.add('hidden');
            });
        });

        function showProductDetails(productName) {
            fetch(`/search?query=${encodeURIComponent(productName)}`)
                .then(response => response.json())
                .then(productDetails => {
                    if (productDetails.error) {
                        document.getElementById('results').innerHTML = `<p>${productDetails.error}</p>`;
                        document.getElementById('product-details').classList.add('hidden');
                        return;
                    }

                    document.getElementById('product-name').textContent = `Name: ${productDetails.name}`;
                    document.getElementById('product-description').textContent = `Description: ${productDetails.description}`;
                    document.getElementById('product-price').textContent = `Price: ${productDetails.base_price}`;
                    document.getElementById('product-manufacturer').textContent = `Manufacturer: ${productDetails.manufacturer}`;
                    document.getElementById('product-details').classList.remove('hidden');
                })
                .catch(error => {
                    console.error('Error fetching product details:', error);
                });
        }
    </script>
</body>
</html>
