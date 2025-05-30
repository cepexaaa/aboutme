#ifndef UNTITLED_BUCKETSTORAGEE_H
#define UNTITLED_BUCKETSTORAGEE_H

#include <cassert>
#include <algorithm>
#include <memory>

template<typename T>
class BucketStorage {
public:

    using difference_type = std::ptrdiff_t; //21:47 добавление пользовательских переменных
    using value_type = T;
    using pointer = T*;
    using reference = T&;
    using iterator_category = std::forward_iterator_tag;
	//size_type, const_reference;

    // Типы итераторов
    class iterator;
    class const_iterator;

    // Конструкторы
    BucketStorage() : BucketStorage(64) {}
    explicit BucketStorage(size_t block_capacity) : blockCapacity(block_capacity), blockCount(1) { //22:50 изменено для массива
        Block blocks = new Block[blockCount];
        currentBlock = &blocks[0];
        currentPos = 0;
    }
    BucketStorage(const BucketStorage& other) : blockCapacity(other.blockCapacity) {
        for (const auto& block : other.blocks) {
            allocateBlock();
            std::copy(block->data, block->data + blockCapacity, currentBlock->data);
        }
    }
    BucketStorage(BucketStorage&& other) noexcept : blockCapacity(other.blockCapacity), blocks(std::move(other.blocks)) { //!!!!ПЕРЕДЕЛАТЬ
        if (!blocks->size) {
            currentBlock = blocks->data[blocks->quanityBlock];//если блок пустой -сносим к херам
            currentPos = other.currentPos;
        } else {
            currentBlock = nullptr;
            currentPos = 0;
        }
    }

    // Деструктор
    ~BucketStorage() { //22:50 изменено для массива
        for (size_t i = 0; i < blockCount; ++i) {
            deallocateBlock(this);//&blocks[i]);
        }
        delete[] this;//blocks;
    }

    // Операторы присваивания
    BucketStorage& operator=(const BucketStorage& other);
    BucketStorage& operator=(BucketStorage&& other) noexcept;

    // Методы
    iterator insert(const value_type& value); //21:47 изменены входные данные
    iterator insert(value_type&& value);
    iterator erase(const_iterator it);
    bool empty() const;
    size_t size() const;
    size_t capacity() const;
    void shrink_to_fit();
    void clear();
    void swap(BucketStorage& other);

    // Итераторы
    iterator begin(); //все ок менять не надо
    const_iterator begin() const;
    const_iterator cbegin() const;
    iterator end();
    const_iterator end() const;
    const_iterator cend() const;

    // Дополнительные методы
    iterator get_to_distance(iterator it, difference_type); //21:51 изменено на difference_type ЕГО НЕ НАДО ЗА 0(1)

private:
    struct Block {
		// massiv_iteratorov[block_cfpfcity]
        T* data; // Заменено на массив //каждый блоr - массив. Есть массив, показывающий активность элементов и ссылки на другие блоки
        size_t size;// в блоке 4 указателя и список на свободные блоки
        size_t capacity;
		int quanityBlock = 0;

        Block(size_t capacity) : data(new T[capacity]), size(0), capacity(capacity) {}
        ~Block() { delete[]data; } // Добавлен деструктор для освобождения памяти
    };

    Block* blocks; // Заменено на массив
    Block* currentBlock;
    size_t currentPos;
    size_t blockCapacity;
    size_t blockCount; // Новое поле для отслеживания количества блоков

    // Вспомогательные методы
    void allocateBlock();
    void deallocateBlock(Block* block);
};

// Реализация итераторов
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
                currentBlock++;
            } while (currentBlock->size == 0);
        }
        return *this;
    }

    iterator operator++(int) {
        iterator tmp = *this;
        ++(*this);
        return tmp;
    }

    friend bool operator==(const iterator& a, const iterator& b) {
        return a.currentBlock == b.currentBlock && a.currentPos == b.currentPos;
    }

    friend bool operator!=(const iterator& a, const iterator& b) {
        return !(a == b);
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
                currentBlock++;
            } while (currentBlock->size == 0);
        }
        return *this;
    }

    const_iterator operator++(int) {
        const_iterator tmp = *this;
        ++(*this);
        return tmp;
    }

    friend bool operator==(const const_iterator& a, const const_iterator& b) {
        return a.currentBlock == b.currentBlock && a.currentPos == b.currentPos;
    }

    friend bool operator!=(const const_iterator& a, const const_iterator& b) {
        return !(a == b);
    }

    const Block* currentBlock;
    size_t currentPos;
};

// Реализация методов
template<typename T>
typename BucketStorage<T>::iterator BucketStorage<T>::insert(const T& value) {
    if (currentBlock->size == blockCapacity) {
        allocateBlock();
    }

    new (&currentBlock->data[currentBlock->size]) T(value);
    currentBlock->size++;

    return iterator(currentBlock, currentBlock->size - 1);
}

template<typename T>
typename BucketStorage<T>::iterator BucketStorage<T>::insert(T&& value) {
    if (currentBlock->size == blockCapacity) {
        allocateBlock();
    }

    new (&currentBlock->data[currentBlock->size]) T(std::move(value));
    currentBlock->size++;

    return iterator(currentBlock, currentBlock->size - 1);
}

template<typename T>
typename BucketStorage<T>::iterator BucketStorage<T>::erase(const_iterator it) {////////////////////////block isn't array "blocks[0]" - is wrong
	assert(it.currentBlock >= blocks[0].get() && it.currentBlock <= blocks.back().get());//blocks->data[blocks->quanityBlock]
    assert(it.currentPos >= 0 && it.currentPos < it.currentBlock->size);

    it.currentBlock->data[it.currentPos].~T();

    for (size_t i = it.currentPos; i < it.currentBlock->size - 1; ++i) {
        new (&it.currentBlock->data[i]) T(std::move(it.currentBlock->data[i + 1]));
        it.currentBlock->data[i + 1].~T();
    }

    it.currentBlock->size--;

    if (it.currentBlock->size == 0 && blocks->quanityBlock > 1) {
        deallocateBlock(it.currentBlock);
        blocks.pop_back();
        currentBlock = blocks.back().get();//blocks->data[blocks->quanityBlock]
    }

    iterator nextIt = it;
    ++nextIt;
    return nextIt;
}

// Реализация методов
template<typename T>
bool BucketStorage<T>::empty() const {
    return blocks->quanityBlock == 0 || (blocks->quanityBlock == 1 && blocks[0]->size == 0);
}

template<typename T>
size_t BucketStorage<T>::size() const {
    if (blocks->quanityBlock == 0) {
        return 0;
    }
    size_t totalSize = 0;
    for (const auto& block : blocks) {
        totalSize += block->size;
    }
    return totalSize;
}

template<typename T>
size_t BucketStorage<T>::capacity() const {
    return blocks->quanityBlock * blockCapacity;
}

template<typename T>
void BucketStorage<T>::shrink_to_fit() {
    // Подсчитываем количество элементов
    size_t totalSize = size();

    // Если количество блоков больше 1, то создаем новый блок и копируем все элементы
    if (blocks->quanityBlock > 1) {
        // Выделяем новый блок с минимальной емкостью, достаточной для хранения всех элементов
        std::unique_ptr<Block> newBlock = std::make_unique<Block>(totalSize);

        // Копируем все элементы в новый блок
        size_t index = 0;
        for (const auto& block : blocks) {
            for (size_t i = 0; i < block->size; ++i) {
                new (&newBlock->data[index]) T(std::move(block->data[i]));
                block->data[i].~T(); // Вызываем деструктор для элементов старых блоков
                ++index;
            }
        }

        // Очищаем старые блоки и сохраняем только новый блок
        blocks.clear();
        blocks.push_back(std::move(newBlock));
        currentBlock = blocks.back().get();//blocks->data[blocks->quanityBlock]
        currentPos = totalSize;
    }

    // Если количество блоков равно 1, то просто обновляем емкость блока
    if (blocks->quanityBlock == 1) {
        std::unique_ptr<Block> newBlock = std::make_unique<Block>(totalSize);
        for (size_t i = 0; i < totalSize; ++i) {
            new (&newBlock->data[i]) T(std::move(blocks[0]->data[i]));
            blocks[0]->data[i].~T(); // Вызываем деструктор для элементов старого блока
        }
        blocks[0] = std::move(newBlock);
        currentBlock = blocks[0].get();
        currentPos = totalSize;
    }
}

template<typename T>
void BucketStorage<T>::clear() {
    blocks.clear();
    currentBlock = nullptr;
    currentPos = 0;
}

template<typename T>
void BucketStorage<T>::swap(BucketStorage& other) {
    std::swap(blocks, other.blocks);
    std::swap(currentBlock, other.currentBlock);
    std::swap(currentPos, other.currentPos);
    std::swap(blockCapacity, other.blockCapacity);
}

template<typename T>
typename BucketStorage<T>::iterator BucketStorage<T>::begin() {
    if (blocks->quanityBlock == 0) {
        return end();
    }
    return iterator(blocks[0].get(), 0);
}

template<typename T>
typename BucketStorage<T>::const_iterator BucketStorage<T>::begin() const {
    if (blocks->quanityBlock == 0) {
        return end();
    }
    return const_iterator(blocks[0].get(), 0);
}

template<typename T>
typename BucketStorage<T>::const_iterator BucketStorage<T>::cbegin() const {
    return begin();
}

template<typename T>
typename BucketStorage<T>::iterator BucketStorage<T>::end() {
    if (blocks->quanityBlock == 0) {
        return iterator(nullptr, 0);
    }
    return iterator(blocks.back().get(), blocks.back()->size);//blocks->data[blocks->quanityBlock]
}

template<typename T>
typename BucketStorage<T>::const_iterator BucketStorage<T>::end() const {
    if (blocks->quanityBlock == 0) {
        return const_iterator(nullptr, 0);
    }
    return const_iterator(blocks.back().get(), blocks.back()->size);//blocks->data[blocks->quanityBlock]
}

template<typename T>
typename BucketStorage<T>::const_iterator BucketStorage<T>::cend() const {
    return end();
}

template<typename T>
typename BucketStorage<T>::iterator BucketStorage<T>::get_to_distance(iterator it, size_t distance) {
    // Если расстояние равно 0, возвращаем текущий итератор
    if (distance == 0) {
        return it;
    }

    // Если расстояние положительное, двигаемся вперед
    if (distance > 0) {
        while (distance > 0) {
            // Вычисляем, сколько элементов можно пройти в текущем блоке
            size_t elementsInBlock = it.currentBlock->size - it.currentPos;
            if (elementsInBlock > distance) {
                // Если расстояние меньше, чем элементов в блоке, просто сдвигаем позицию
                it.currentPos += distance;
                distance = 0;
            } else {
                // Иначе, переходим к следующему блоку
                distance -= elementsInBlock;
                do {
                    ++it.currentBlock;
                } while (it.currentBlock->size == 0);
                it.currentPos = 0;
            }
        }
    } else {
        // Если расстояние отрицательное, двигаемся назад
        while (distance < 0) {
            // Вычисляем, сколько элементов можно пройти назад в текущем блоке
            size_t elementsInBlock = it.currentPos;
            if (elementsInBlock >= -distance) {
                // Если расстояние меньше, чем элементов в блоке, просто сдвигаем позицию
                it.currentPos += distance;
                distance = 0;
            } else {
                // Иначе, переходим к предыдущему блоку
                distance += elementsInBlock;
                do {
                    --it.currentBlock;
                } while (it.currentBlock->size == 0);
                it.currentPos = it.currentBlock->size - 1;
            }
        }
    }

    return it;
}

// Вспомогательные методы
template<typename T>
void BucketStorage<T>::allocateBlock() { //22:49
    Block* newBlocks = new Block[blockCount + 1];
    std::copy(blocks, blocks + blockCount, newBlocks);
    delete[] blocks;
    blocks = newBlocks;
    currentBlock = &blocks[blockCount];
    currentPos = 0;
    blockCount++;
}

template<typename T>
void BucketStorage<T>::deallocateBlock(Block* block) { //22:49
    // Освобождаем память блока
    delete[] block->data;
    block->data = nullptr;
    block->size = 0;
    block->capacity = 0;
}


// Операторы присваивания
template<typename T>
BucketStorage<T>& BucketStorage<T>::operator=(const BucketStorage& other) {
    if (this != &other) {
        clear();
        blockCapacity = other.blockCapacity;
        for (const auto& block : other.blocks) {
            allocateBlock();
            std::copy(block->data.get(), block->data.get() + block->size, currentBlock->data.get());
            currentBlock->size = block->size;
        }
    }
    return *this;
}

template<typename T>
BucketStorage<T>& BucketStorage<T>::operator=(BucketStorage&& other) noexcept {
    if (this != &other) {
        clear();
        blockCapacity = other.blockCapacity;
        blocks = std::move(other.blocks);
        if (blocks->quanityBlock != 0) {
            currentBlock = blocks.back().get();//blocks->data[blocks->quanityBlock]
            currentPos = other.currentPos;
        } else {
            currentBlock = nullptr;
            currentPos = 0;
        }
    }
    return *this;
}
#endif
