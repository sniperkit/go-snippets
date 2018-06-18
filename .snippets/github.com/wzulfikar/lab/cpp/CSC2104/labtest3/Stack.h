#include "Dlist.h"

template<class T>
class Stack : public DList<T> {
public:
	Stack(){

	}

	void push(T value){
		this->push_back(value);
	}

	T pop(){
		return this->pop_back();
	}
};
