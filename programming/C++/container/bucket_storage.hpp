#ifndef BUCKET_STORAGE_HPP
#define BUCKET_STORAGE_HPP

#include <vector>
#include <cassert>
#include <algorithm>
#include <memory>

template<typename T>
class BucketStorage {
  public:
	class iterator;
	class const_iterator;

	using value_type = T;
	using reference = T&;
	using const_reference = const T&;
//	using iterator = iterator<T>;
//	using const_iterator = const iterator<T>;
	using difference_type = std::ptrdiff_t;
	using size_type = size_t;


	// Конструкторы
	BucketStorage() : BucketStorage(64) {}
	explicit BucketStorage(size_t block_capacity) : currentBlock(nullptr), currentPos(0), freeBlocks(0), blockCapacity(block_capacity) {}
	BucketStorage(const BucketStorage& other) : blockCapacity(other.blockCapacity) {
		for (const auto& block : other.blocks) {
			allocateBlock();
			memcpy(currentBlock->data, block->data, blockCapacity * sizeof(T));
			currentPos = other.currentPos;
		}
	}
	BucketStorage(BucketStorage&& other) noexcept : blockCapacity(other.blockCapacity), blocks(std::move(other.blocks)), currentBlock(other.currentBlock), currentPos(other.currentPos), freeBlocks(other.freeBlocks) {
		other.currentBlock = nullptr;
		other.currentPos = 0;
		other.freeBlocks = 0;
	}

	// Деструктор
	~BucketStorage() { clear(); }

	// Операторы присваивания
	BucketStorage& operator=(const BucketStorage& other) {
		if (this != &other) {
			clear();
			blockCapacity = other.blockCapacity;
			for (const auto& block : other.blocks) {
				allocateBlock();
				memcpy(currentBlock->data, block->data, blockCapacity * sizeof(T));
				currentPos = other.currentPos;
			}
		}
		return *this;
	}
	BucketStorage& operator=(BucketStorage&& other) noexcept {
		if (this != &other) {
			clear();
			blockCapacity = other.blockCapacity;
			blocks = std::move(other.blocks);
			currentBlock = other.currentBlock;
			currentPos = other.currentPos;
			freeBlocks = other.freeBlocks;
			other.currentBlock = nullptr;
			other.currentPos = 0;
			other.freeBlocks = 0;
		}
		return *this;
	}

    //влепил кучу noexcept - правильно ли было так сделать?
	iterator insert(const value_type& value);
	iterator insert(value_type&& value);
	iterator erase(const_iterator it);
	iterator erase(iterator it);//пришлось дублировать, так как нужен и такой итератор и другой
	bool empty() const noexcept;
	size_t size() const noexcept;
	size_t capacity() const noexcept;
	void shrink_to_fit();
	void clear();
	void swap(BucketStorage& other) noexcept;

	// Итераторы
	iterator begin() noexcept;
	const_iterator begin() const noexcept;
	const_iterator cbegin() const noexcept;
	iterator end() noexcept;
	const_iterator end() const noexcept;
	const_iterator cend() const noexcept;

	// Дополнительные методы
	iterator get_to_distance(iterator it, const difference_type distance);

  private:
    size_t blockCapacity; // как сделать так, чтобы не передавать в блок блокКапасити
	struct Block {
		T* data;//T data[blockCapacity]; // инициализировать массив типа Т на капасити элементов
        size_t size;
        size_t blockCapacity;
        int* identificators;
		//int nextActive;// - ?
        Block* nextBlock;
        Block* prevBlock;

        Block(size_t blockCapacity) : data(new T[blockCapacity]), size(0), blockCapacity(blockCapacity), identificators(new int[blockCapacity]) {}
        ~Block() { delete[]data; }//доработать диструктор
	};

	std::vector<Block*> blocks;
	Block* currentBlock;
	int currentPos;
	int freeBlocks;
	 // Вместимость блока

	// Вспомогательные методы
	void allocateBlock();
	void deallocateBlock(Block* block);
};



template<typename T>
class BucketStorage<T>::iterator {
public:
    using difference_type = std::ptrdiff_t;
    using value_type = T;
    using pointer = T*;
    using reference = T&;
    using iterator_category = std::forward_iterator_tag;

    iterator(Block* block, size_t pos) : currentBlock(block), currentPos(pos) {}

    reference operator*() const { return currentBlock->data[currentPos]; }
    pointer operator->() const { return &currentBlock->data[currentPos]; }

    iterator& operator++() {
        currentPos++;
        if (currentPos == currentBlock->size) {
            currentPos = 0;
            do {
                currentBlock = currentBlock->nextBlock;
            } while (currentBlock && currentBlock->size == 0);
        }
        return *this;
    }

    iterator operator++(int) {
        iterator tmp = *this;
        ++(*this);
        return tmp;
    }

    iterator operator--() {
        return *this;
    }

    iterator operator--(int) {
        return *this;
    }

    friend bool operator==(iterator& a, iterator& b) {
        return a.currentBlock == b.currentBlock && a.currentPos == b.currentPos;
    }

    friend bool operator!=(const iterator& a, const iterator& b) {
        return !(a == b);
    }

    friend bool operator!=(const iterator& a, const const_iterator& b) {
        return !(a.currentBlock == b.currentBlock && a.currentPos == b.currentPos);
    }

    friend bool operator<(const iterator& a, const iterator& b) {
        return true;
    }
    friend bool operator>(const iterator& a, const iterator& b) {
        return true;
    }
    friend bool operator<=(const iterator& a, const iterator& b) {
        return a == b || a < b;
    }
    friend bool operator>=(const iterator& a, const iterator& b) {
        return a == b || a > b;
    }

    Block* currentBlock;
    size_t currentPos;
};

template<typename T>
class BucketStorage<T>::const_iterator {
public:
    using difference_type = std::ptrdiff_t;
    using value_type = const T;
    using pointer = const T*;
    using reference = const T&;
    using iterator_category = std::forward_iterator_tag;

    const_iterator(const Block* block, size_t pos) : currentBlock(block), currentPos(pos) {}

    reference operator*() const { return currentBlock->data[currentPos]; }
    pointer operator->() const { return &currentBlock->data[currentPos]; }

    const_iterator& operator++() {
        currentPos++;
        if (currentPos == currentBlock->size) {
            currentPos = 0;
            do {
                currentBlock = currentBlock->nextBlock;
            } while (currentBlock && currentBlock->size == 0);
        }
        return *this;
    }

    const_iterator operator++(int) {
        const_iterator tmp = *this;
        ++(*this);
        return tmp;
    }

    friend bool operator==(const_iterator& a, const_iterator& b) {
        return a.currentBlock == b.currentBlock && a.currentPos == b.currentPos;
    }

    friend bool operator!=(const const_iterator& a, const const_iterator& b) {
        return !(a == b);
    }

    friend bool operator!=(const const_iterator& a, const iterator& b) {
        return !(a.currentBlock == b.currentBlock && a.currentPos == b.currentPos);
    }


//    friend bool operator!=(const iterator& a, const const_iterator& b) {
//        return !(a.currentBlock == b.currentBlock && a.currentPos == b.currentPos);
//    }

//    friend bool operator!=(const const_iterator& a, const iterator& b) {
//        return !(a.currentBlock == b.currentBlock && a.currentPos == b.currentPos);
//    }

    const_iterator operator--() {
        return *this;
    }

    const_iterator operator--(int) {
        return *this;
    }

    friend bool operator<(const const_iterator& a, const const_iterator& b) {
        return true;
    }
    friend bool operator>(const const_iterator& a, const const_iterator& b) {
        return true;
    }
    friend bool operator<=(const const_iterator& a, const const_iterator& b) {
        return a == b || a < b;
    }
    friend bool operator>=(const const_iterator& a, const const_iterator& b) {
        return a == b || a > b;
    }

    const Block* currentBlock;
    size_t currentPos;
};





template<typename T>
typename BucketStorage<T>::iterator BucketStorage<T>::insert(const value_type& value) {
	// Реализация вставки с использованием copy-constructor
    iterator it  = new iterator();
    return it;
}

template<typename T>
typename BucketStorage<T>::iterator BucketStorage<T>::insert(value_type&& value) {
	// Реализация вставки с использованием move-constructor
    iterator it = new iterator();
    return it;
}

template<typename T>
typename BucketStorage<T>::iterator BucketStorage<T>::erase(const_iterator it) {
	// Реализация удаления элемента
    return it;
}

template<typename T>
typename BucketStorage<T>::iterator BucketStorage<T>::erase(iterator it) {
    // Реализация удаления элемента
    return it;
}

template<typename T>
bool BucketStorage<T>::empty() const noexcept{
	// Проверка, является ли контейнер пустым
    return false;
}

template<typename T>
size_t BucketStorage<T>::size() const noexcept{
	// Возвращение количества элементов в контейнере
    return 0;
}

template<typename T>
size_t BucketStorage<T>::capacity() const noexcept{
	// Возвращение максимального количества элементов, которые могут быть сохранены без расширения
    return 0;
}

template<typename T>
void BucketStorage<T>::shrink_to_fit() {
	// Изменение ёмкости контейнера до минимально необходимого
}

template<typename T>
void BucketStorage<T>::clear() {
	// Очистка контейнера
}

template<typename T>
void BucketStorage<T>::swap(BucketStorage& other) noexcept {
	// Обмен содержимым контейнеров
}

// Реализация итераторов
// ...

template<typename T>
typename BucketStorage<T>::iterator BucketStorage<T>::get_to_distance(iterator it, const difference_type distance) {
	// Сдвиг итератора на указанное расстояние
    return it;
}

//template<typename T>
//typename BucketStorage<T>::const_iterator BucketStorage<T>::begin() const noexcept{
//    return nullptr;
//}









#endif //BUCKET_STORAGE_HPP
