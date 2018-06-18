#include "Dlist.h"

template<class T>
class Queue : public DList<T> {
public:
	Queue () {

	}

	void push(T value){
		this->push_back(value);
	}

	T pop(){
		return this->pop_front();
	}
};
