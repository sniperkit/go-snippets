// This file is part of the Orbbec Astra SDK [https://orbbec3d.com]
// Copyright (c) 2015 Orbbec 3D
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// Be excellent to each other.
#ifndef ASTRA_SIGNAL_H
#define ASTRA_SIGNAL_H

#include <functional>

//TODO: make add/remove callbacks threadsafe with mutex

namespace astra {

    template<typename R, typename... Signature> class CallbackList;

    template<typename R, typename... Signature>
    class CallbackNode
    {

    public:

        typedef std::function<R (Signature...)> callback_type;

        CallbackNode(const callback_type& cb)
            : prev_(nullptr), next_(nullptr), callback_(cb), refCount_(1)
        {}

        void inc_reference() { refCount_++; }
        void dec_reference()
        {
            refCount_--;
            if (refCount_ == 0)
            {
                delete this;
            }
        }

        void deactivate()
        {
            callback_ = nullptr;
        }

        void remove_self()
        {
            deactivate();
            if (next_)
            {
                next_->prev_ = prev_;
            }
            if (prev_)
            {
                prev_->next_ = next_;
            }
            dec_reference();
        }

    private:

        CallbackNode* prev_;
        CallbackNode* next_;
        callback_type callback_;
        unsigned int refCount_;

        friend class CallbackList<R, Signature...>;
    };

    template<typename R, typename... Signature>
    class CallbackList
    {
    public:

        typedef CallbackNode<R, Signature...> node_type;
        typedef typename node_type::callback_type callback_type;

        CallbackList()
            : count_(0), head_(nullptr) { };

        ~CallbackList()
        {
            if (!head_)
                return;

            node_type* node = head_;
            while(node->next_ != head_)
            {
                node->next_->remove_self();
            }
            head_->dec_reference();
            head_->dec_reference();
        }

        size_t add(const callback_type& cb)
        {
            node_type* node = new node_type(cb);

            if (!head_)
            {
                head_ = node;
                head_->inc_reference();
                // set up as ring list
                head_->next_ = head_;
                head_->prev_ = head_;
            }
            else
            {
                node->prev_ = head_->prev_;
                node->next_ = head_;

                if (head_->prev_)
                {
                    head_->prev_->next_ = node;
                }

                head_->prev_ = node;
            }

            count_++;

            return size_t(node);
        }

        bool remove(size_t id)
        {
            if (!head_)
                return false;

            if (size_t(head_) == id)
            {
                head_->deactivate();
                count_--;
                return true;
            }

            for (node_type* node = head_->next_ ? head_->next_ : head_;
                 node != head_;
                 node = node->next_)
            {
                if (id == size_t(node))
                {
                    node->remove_self();
                    count_--;
                    return true;
                }
            }

            return false;
        }

        unsigned debug_count()
        {
            if (!head_)
                return 0;

            unsigned count = 0;
            for (node_type* node = head_->next_ ? head_->next_ : head_;
                 node != head_;
                 node = node->next_)
            {
                count++;
            }

            count++;

            return count;
        }
        unsigned int& count()
        {
            return count_;
        }

        void invoke(Signature... sig)
        {
            if (!head_)
                return;

            node_type* node = head_;
            node->inc_reference();
            do
            {
                if (node->callback_ != nullptr)
                {
                    node->callback_(sig...);
                }
                node_type* prev = node;
                node = node->next_;
                node->inc_reference();
                prev->dec_reference();
            } while (node != head_);
            node->dec_reference();
        }

    private:

        unsigned int count_;
        node_type* head_;
    };

    template<typename... Args>
    class signal
    {
        typedef CallbackList<void, Args...> callback_list_t;

    public:

        typedef typename callback_list_t::callback_type callback_type;

        signal()
            : callbackList_() { }
        unsigned int& slot_count()
        {
            return callbackList_.count();
        }
        size_t operator+=(const callback_type& cb)
        {
            return callbackList_.add(cb);
        }
        bool operator-=(size_t id)
        {
            return callbackList_.remove(id);
        }

        void raise(Args... args)
        {
            callbackList_.invoke(args...);
        }
    private:
        callback_list_t callbackList_;

    };

    template<>
    class signal<void>
    {
        typedef CallbackList<void, bool> callback_list_t;

    public:
        typedef std::function<void ()> callback_type;

        signal()
            : callbackList_() { }
        unsigned int& slot_count()
        {
            return callbackList_.count();
        }
        size_t operator+=(const callback_type& cb)
        {
            return callbackList_.add([cb] (bool b) { cb(); });
        }

        bool operator-=(size_t id)
        {
            return callbackList_.remove(id);
        }

        void raise()
        {
            callbackList_.invoke(true);
        }
    private:
        callback_list_t callbackList_;

    };
}

#endif /* ASTRA_SIGNAL_H */
