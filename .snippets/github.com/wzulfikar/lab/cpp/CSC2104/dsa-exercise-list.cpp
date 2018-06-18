#include <iostream>
#include <list>

using namespace std;

class Product
{
  public:
    string name;
    Product(string name) {
        this->name = name;
    }
};

typedef list<Product*> ProductList;
ProductList::iterator it;

void displayProducts (ProductList products) {
	int counter = 1;

	printf("\nDisplaying product(s)..\n");
	for(it = products.begin(); it != products.end(); it++)
		printf("  %d. %s\n", counter++, (*it)->name.c_str());

	printf("\n%lu product(s) succesfully displayed ✔\n", products.size());
}

void appendProduct (Product* p, ProductList *products) {
	products->push_back(p);
	printf("Product added: %s\n", p->name.c_str());
}

void deleteProduct (string productName, ProductList *products) {
	bool isDeleted = false;
	for(it = products->begin(); it != products->end(); it++) {
		if ((*it)->name == productName) {
			isDeleted = true;
			products->remove(*it);
			printf("\n✔ Product '%s' has been succesfully deleted.\n", 
				productName.c_str()
			);
			break;
		}
	}

	if (!isDeleted) {
		printf("\n❗️ Failed to delete product:\n  Product '%s' doesn't exist.\n", 
			productName.c_str()
		);
	}
}

void insertProductAtIndex (Product *product, 
									int insertAtIndex, 
									ProductList *products) {

	const unsigned long MAX_INDEX = products->size() - 1;
	if (insertAtIndex > MAX_INDEX) {
		printf("\n❗️ Failed to insert product '%s' at index %d:\n  Index shouldn't exceed %lu\n", 
			product->name.c_str(), 
			insertAtIndex,
			MAX_INDEX
		);
		return;
	}

	int counter = 0;
	for(it = products->begin(); it != products->end(); it++) {
		if (counter == insertAtIndex) {
			products->insert(it, product);
			break;
		}
		counter++;
	}
}

// TASKS:
// - create list ✔
// - append new product ✔
// - insert product to specific location ✔
// - delete product ✔
// - display ✔
int main()
{
	ProductList	products;

	appendProduct(new Product("yoloo"), &products);
	appendProduct(new Product("second product"), &products);
	appendProduct(new Product("third product"), &products);
	appendProduct(new Product("forth product"), &products);

	// deleteProduct("yoloo", &products);

	insertProductAtIndex(new Product("fifth product at idx 2"), 2, &products);
	
	displayProducts(products);

	return 0;
}
