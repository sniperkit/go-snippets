#ifndef DLIST_H
#define DLIST_H

#include <iostream>

using namespace std;

template<class T>
class Item{
	public:
	T value;
	Item<T>* next;
	Item<T>* prev;
};

template<class T>
class DList{
	protected:
		Item<T>* head;
		Item<T>* tail;
	public:
		DList(){
			head = new Item<T>();
			tail = new Item<T>();

			head->next = tail;
			head->prev = NULL;
			
			tail->next = NULL;
			tail->prev = head;
		}

		virtual void push(T value)=0;
		virtual T pop()=0;

		void push_back(T value){
			Item<T>* item = new Item<T>();
			item->value = value;
			item->next = tail;
			item->prev = tail->prev;

			tail->prev->next = item;
			tail->prev = item;
		}

		T pop_back(){
			if(isEmpty()){
				cout<<"List is empty, nothing to pop.";
				return false;
			}else{
				Item<T>* temp = tail->prev;
				tail->prev = temp->prev;
				temp->prev->next = tail;

				T value = temp->value;
				delete temp;
				return value;
			}
		}

		T back(){
			if(!isEmpty()){
				return tail->prev->value;
			}
		}

		void push_front(T value){
			Item<T>* item = new Item<T>();
			item->value = value;
			item->next = head->next;
			item->prev = head;

			head->next->prev = item;
			head->next = item;
		}

		T pop_front(){
			if(isEmpty()){
				cout<<"List is empty, nothing to pop.";
			}else{
				Item<T>* temp = head->next;
				head->next = temp->next;
				temp->next->prev = head;
				
				T value = temp->value;
				delete temp;
				return value;
			}
		}
		
		T front(){
			if(!isEmpty()){
				return head->next->value;
			}
		}

		bool isEmpty(){
			return head->next==tail;
		}

		void display(bool reverse=false){
			Item<T>* temp = new Item<T>();
			if(reverse){
				temp = tail;
				while(temp->prev!=head){
					cout<<temp->prev->value<<" -> ";
					temp = temp->prev;
				}
			}else{
				temp = head;
				while(temp->next!=tail){
					cout<<temp->next->value<<" -> ";
					temp = temp->next;
				}
			}
			cout<<"END"<<endl;
		}
};

#endif
