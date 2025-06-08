#ifndef BUCKET_STORAGE_HPP
#define BUCKET_STORAGE_HPP

#include <memory>

template< typename T >
class BucketStorage
{
	template< bool isConst >
	class BucketStorageIterator;
	template< bool isConst >
	friend class BucketStorageIterator;
	template< bool isConstRight >
	friend class BucketStorageIterator;

  public:
	using iterator = BucketStorageIterator< false >;
	using const_iterator = BucketStorageIterator< true >;
	using value_type = T;
	using reference = T&;
	using const_reference = const T&;
	using difference_type = std::ptrdiff_t;
	using size_type = std::size_t;

	BucketStorage(const size_type block_capacity = 64);
	BucketStorage(const BucketStorage& other);
	BucketStorage(BucketStorage&& other) noexcept;
	~BucketStorage();

	BucketStorage& operator=(BucketStorage&& other) noexcept;
	BucketStorage& operator=(const BucketStorage& other);

	iterator insert(const value_type& value);
	iterator insert(value_type&& value);
	iterator erase(const_iterator it);
	bool empty() const noexcept;
	size_type size() const noexcept;
	size_type capacity() const noexcept;
	void shrink_to_fit();
	void clear();
	void swap(BucketStorage& other) noexcept;
	iterator begin() noexcept;
	const_iterator begin() const noexcept;
	const_iterator cbegin() const noexcept;
	iterator end() noexcept;
	const_iterator end() const noexcept;
	const_iterator cend() const noexcept;
	iterator get_to_distance(iterator it, const difference_type distance);

  private:
	struct Block;
	size_type blockCapacity = 64;
	int c_countBlocks = 0;
	size_type c_bucketSize = 0;
	Block* b_firstBlock = nullptr;
	Block* b_lastBlock = nullptr;
	Block* b_firstEmptyBlock = nullptr;
	Block* b_lastEmptyBlock = b_firstEmptyBlock;
	int freeBlocks = 0;

	struct Block
	{
		size_type b_size;
		size_type b_id;
		bool b_isFullBlock = false;
		Block* b_nextBlock = nullptr;
		Block* b_prevBlock = nullptr;
		Block* b_nextFreeBlock = nullptr;
		Block* b_prevFreeBlock = nullptr;
		struct Node;
		Node* blockStart;
		int n_countFilled = 0;
		Node* n_thisBegin;
		Node* n_thisEnd;
		Node* n_firstEmpty;
		Node* n_deleteNodeS;

		void deleteNodesB2E(Node* node);
		void deleteNodesB2F(Node* node);

		Block(size_type blockCapacity, value_type newValue, size_type b_id);
		~Block();

		struct Node
		{
			Node(size_type n_id);
			value_type* n_value = static_cast< value_type* >(operator new(sizeof(value_type)));
			size_type n_id;
			Node* n_next = nullptr;
			Node* n_prev = nullptr;
			~Node();
		};
	};

	void lincNodes(BucketStorage::Block::Node* a, BucketStorage::Block::Node* b);
	void simpleErase(const_iterator it);
	template< typename V >
	typename BucketStorage< T >::iterator insertSimple(V&& value);
};

template< typename T >
template< bool isConst >
class BucketStorage< T >::BucketStorageIterator
{
  public:
	using iterator_category = std::bidirectional_iterator_tag;
	using difference_type = std::ptrdiff_t;
	using value_type = T;
	using pointer = std::conditional_t< isConst, const T*, T* >;
	using reference = std::conditional_t< isConst, const T&, T& >;

	BucketStorageIterator(typename BucketStorage< T >::Block* block = nullptr, typename BucketStorage< T >::Block::Node* node = nullptr);
	reference operator*() const;
	pointer operator->() const;
	BucketStorageIterator& operator++();
	BucketStorageIterator operator++(int);
	BucketStorageIterator& operator--();
	BucketStorageIterator operator--(int);
	BucketStorageIterator& operator=(const BucketStorageIterator< !isConst >& other);
	template< bool isConstRight >
	bool operator==(const BucketStorageIterator< isConstRight >& other) const;
	template< bool isConstRight >
	bool operator!=(const BucketStorageIterator< isConstRight >& other) const;
	bool operator<(const BucketStorageIterator& other) const;
	bool operator>(const BucketStorageIterator& other) const;
	bool operator<=(const BucketStorageIterator& other) const;
	bool operator>=(const BucketStorageIterator& other) const;
	operator BucketStorageIterator< !isConst >() const;

  private:
	friend class BucketStorage;
	typename BucketStorage< T >::Block* currentBlock;
	typename BucketStorage< T >::Block::Node* currentNode;

	void createEndNode();
};

template< typename T >
BucketStorage< T >::BucketStorage(const BucketStorage& other) :
	blockCapacity(other.blockCapacity), freeBlocks(other.freeBlocks)
{
	Block* otherBlock = other.b_firstBlock;
	Block* prevBlock = nullptr;
	Block* prevFreeBlock = nullptr;
	b_firstBlock = other.b_firstBlock;
	b_lastBlock = other.b_lastBlock;
	c_countBlocks = other.c_countBlocks;
	c_bucketSize = other.c_bucketSize;
	size_type b_new_id = 0;
	while (otherBlock != nullptr)
	{
		auto newBlock = new Block(otherBlock->b_size, *(otherBlock->n_thisBegin->n_value), b_new_id++);
		newBlock->n_countFilled = otherBlock->n_countFilled;
		newBlock->b_isFullBlock = otherBlock->b_isFullBlock;

		auto* otherNode = otherBlock->n_thisBegin;
		auto* newNode = newBlock->n_thisBegin;
		int i = 0;
		while (otherNode != nullptr && i++ < blockCapacity)
		{
			new (*(&newNode->n_value)) value_type(*(otherNode->n_value));
			newNode->n_next = (otherNode->n_next != nullptr) ? newNode + 1 : nullptr;
			newNode->n_prev = (otherNode->n_prev != nullptr) ? newNode - 1 : nullptr;
			otherNode = otherNode->n_next;
			newNode = newNode->n_next;
			if (otherNode == otherBlock->n_firstEmpty)
			{
				newBlock->n_firstEmpty = newNode;
			}
		}
		newBlock->b_prevBlock = prevBlock;
		if (prevBlock != nullptr)
		{
			prevBlock->b_nextBlock = newBlock;
		}
		else
		{
			b_firstBlock = newBlock;
		}

		if (!otherBlock->b_isFullBlock)
		{
			newBlock->b_prevFreeBlock = prevFreeBlock;
			if (prevFreeBlock != nullptr)
			{
				prevFreeBlock->b_nextFreeBlock = newBlock;
			}
			prevFreeBlock = newBlock;
		}

		if (otherBlock->n_firstEmpty == nullptr)
		{
			newBlock->n_firstEmpty = nullptr;
		}
		prevBlock = newBlock;
		otherBlock = otherBlock->b_nextBlock;
	}

	b_firstEmptyBlock = other.b_firstEmptyBlock;
	b_lastEmptyBlock = other.b_lastEmptyBlock;
}

template< typename T >
BucketStorage< T >::BucketStorage(const size_type block_capacity) : blockCapacity(block_capacity)
{
	if (block_capacity < 1)
	{
		throw std::invalid_argument("Block capacity must be greater than zero\n");
	}
}

template< typename T >
BucketStorage< T >::BucketStorage(BucketStorage&& other) noexcept
{
	this->swap(other);
}

template< typename T >
BucketStorage< T >::~BucketStorage()
{
	clear();
}

template< typename T >
BucketStorage< T >& BucketStorage< T >::operator=(BucketStorage&& other) noexcept
{
	if (&other != this)
	{
		BucketStorage< T > temp(std::move(other));
		this->swap(temp);
	}
	return *this;
}

template< typename T >
BucketStorage< T >& BucketStorage< T >::operator=(const BucketStorage& other)
{
	if (&other != this)
	{
		BucketStorage< T > temp(other);
		this->swap(temp);
	}
	return *this;
}

template< typename T >
BucketStorage< T >::Block::Block(size_type blockCapacity, value_type newValue, size_type b_id) :
	b_size(blockCapacity), b_id(b_id)
{
	Node* nodes = static_cast< Node* >(operator new(blockCapacity * sizeof(Node)));
	blockStart = nodes;

	Node* prevNode = nullptr;
	for (size_type i = 0; i < blockCapacity; ++i)
	{
		new (&nodes[i]) Node(i);
		Node* newNode = &nodes[i];
		newNode->n_prev = prevNode;
		newNode->n_next = nullptr;
		if (prevNode != nullptr)
		{
			prevNode->n_next = newNode;
		}
		prevNode = newNode;
	}

	n_thisBegin = &nodes[0];
	n_thisEnd = &nodes[blockCapacity - 1];
	n_firstEmpty = &nodes[1];
	n_deleteNodeS = n_thisEnd;
}

template< typename T >
BucketStorage< T >::Block::Node::Node(size_type n_id) : n_id(n_id)
{
}

template< typename T >
BucketStorage< T >::Block::Node::~Node()
{
	delete n_value;
}

template< typename T >
BucketStorage< T >::Block::~Block()
{
	if (b_nextBlock == nullptr && n_thisEnd->n_next != nullptr)
	{
		delete n_thisEnd->n_next;
	}
	Node* node = n_thisBegin;
	if (n_firstEmpty == nullptr)
	{
		deleteNodesB2E(node);
	}
	else if (n_deleteNodeS == n_thisEnd)
	{
		if (n_firstEmpty != n_thisBegin)
		{
			deleteNodesB2F(node);
		}
		else
		{
			deleteNodesB2E(node);
		}
	}
	else
	{
		deleteNodesB2F(node);
		node = n_deleteNodeS->n_next;
		deleteNodesB2E(node);
	}
	::operator delete(blockStart);
}

template< typename T >
void BucketStorage< T >::Block::deleteNodesB2E(Node* node)
{
	while (node != nullptr && node != n_thisEnd)
	{
		delete node->n_value;
		node = node->n_next;
	}
	if (n_thisEnd != nullptr && node->n_value != nullptr)
	{
		delete node->n_value;
	}
}
template< typename T >
void BucketStorage< T >::Block::deleteNodesB2F(Node* node)
{
	while (node != nullptr && node != n_firstEmpty)
	{
		delete node->n_value;
		node = node->n_next;
	}
}

template< typename T >
template< bool isConst >
BucketStorage< T >::BucketStorageIterator< isConst >::BucketStorageIterator(
	typename BucketStorage< T >::Block* block,
	typename BucketStorage< T >::Block::Node* node) : currentBlock(block), currentNode(node)
{
}

template< typename T >
template< bool isConst >
typename BucketStorage< T >::template BucketStorageIterator< isConst >::reference
	BucketStorage< T >::BucketStorageIterator< isConst >::operator*() const
{
	if (currentNode == nullptr || (currentBlock->n_thisEnd == currentNode->n_prev))
	{
		throw std::logic_error("Getting the pointer dereference to the default constructable iterator\n");
	}
	return *currentNode->n_value;
}

template< typename T >
template< bool isConst >
typename BucketStorage< T >::template BucketStorageIterator< isConst >::pointer
	BucketStorage< T >::BucketStorageIterator< isConst >::operator->() const
{
	return currentNode->n_value;
}

template< typename T >
template< bool isConst >
BucketStorage< T >::BucketStorageIterator< isConst >& BucketStorage< T >::BucketStorageIterator< isConst >::operator++()
{
	if (currentNode->n_prev == currentBlock->n_thisEnd || (currentBlock->n_firstEmpty && currentNode == currentBlock->n_firstEmpty))
	{
		return *this;
	}
	if (currentNode->n_next)
	{
		if (currentNode->n_next == currentBlock->n_firstEmpty || currentNode == currentBlock->n_thisEnd)
		{
			if (currentBlock->b_nextBlock != nullptr)
			{
				currentBlock = currentBlock->b_nextBlock;
				currentNode = currentBlock->n_thisBegin;
			}
			else
			{
				if (currentNode->n_next == currentBlock->n_firstEmpty)
				{
					currentNode = currentNode->n_next;
				}
				else
				{
					createEndNode();
				}
			}
		}
		else
		{
			currentNode = currentNode->n_next;
		}
	}
	else
	{
		createEndNode();
	}
	return *this;
}

template< typename T >
template< bool isConst >
BucketStorage< T >::BucketStorageIterator< isConst > BucketStorage< T >::BucketStorageIterator< isConst >::operator++(int)
{
	BucketStorageIterator tmp = *this;
	++(*this);
	return tmp;
}

template< typename T >
template< bool isConst >
BucketStorage< T >::BucketStorageIterator< isConst >& BucketStorage< T >::BucketStorageIterator< isConst >::operator--()
{
	if (currentNode->n_prev && currentNode == currentBlock->n_thisBegin)
	{
		if (currentBlock->b_prevBlock->n_firstEmpty)
		{
			currentNode = currentBlock->b_prevBlock->n_firstEmpty->n_prev;
		}
		else
		{
			currentNode = currentBlock->b_prevBlock->n_thisEnd;
		}
		currentBlock = currentBlock->b_prevBlock;
	}
	else if (currentNode->n_prev)
	{
		currentNode = currentNode->n_prev;
	}
	return *this;
}

template< typename T >
template< bool isConst >
BucketStorage< T >::BucketStorageIterator< isConst > BucketStorage< T >::BucketStorageIterator< isConst >::operator--(int)
{
	BucketStorageIterator tmp = *this;
	--(*this);
	return tmp;
}

template< typename T >
template< bool isConst >
BucketStorage< T >::BucketStorageIterator< isConst >&
	BucketStorage< T >::BucketStorageIterator< isConst >::operator=(const BucketStorageIterator< !isConst >& other)
{
	currentBlock = other.currentBlock;
	currentNode = other.currentNode;
	return *this;
}

template< typename T >
template< bool isConst >
template< bool isConstRight >
bool BucketStorage< T >::BucketStorageIterator< isConst >::operator==(const BucketStorage< T >::BucketStorageIterator< isConstRight >& other) const
{
	return currentNode == other.currentNode;
}

template< typename T >
template< bool isConst >
template< bool isConstRight >
bool BucketStorage< T >::BucketStorageIterator< isConst >::operator!=(const BucketStorage< T >::BucketStorageIterator< isConstRight >& other) const
{
	return !(*this == other);
}

template< typename T >
template< bool isConst >
bool BucketStorage< T >::BucketStorageIterator< isConst >::operator<(const BucketStorageIterator& other) const
{
	if (currentBlock != other.currentBlock)
	{
		return currentBlock->b_id < other.currentBlock->b_id;
	}
	else
	{
		return currentNode->n_id < other.currentNode->n_id;
	}
}

template< typename T >
template< bool isConst >
bool BucketStorage< T >::BucketStorageIterator< isConst >::operator>(const BucketStorageIterator& other) const
{
	if (currentBlock != other.currentBlock)
	{
		return currentBlock->b_id > other.currentBlock->b_id;
	}
	else
	{
		return currentNode->n_id > other.currentNode->n_id;
	}
}

template< typename T >
template< bool isConst >
bool BucketStorage< T >::BucketStorageIterator< isConst >::operator<=(const BucketStorageIterator& other) const
{
	return *this == other || *this < other;
}

template< typename T >
template< bool isConst >
bool BucketStorage< T >::BucketStorageIterator< isConst >::operator>=(const BucketStorageIterator& other) const
{
	return *this == other || *this > other;
}

template< typename T >
template< bool isConst >
BucketStorage< T >::BucketStorageIterator< isConst >::operator BucketStorageIterator< !isConst >() const
{
	return BucketStorageIterator< !isConst >(const_cast< Block* >(currentBlock), const_cast< Block::Node* >(currentNode));
}

template< typename T >
template< bool isConst >
void BucketStorage< T >::BucketStorageIterator< isConst >::createEndNode()
{
	auto endNode = new typename BucketStorage< T >::Block::Node(currentNode->n_id + 1);
	currentNode->n_next = endNode;
	endNode->n_prev = currentNode;
	currentNode = endNode;
}

template< typename T >
void BucketStorage< T >::lincNodes(BucketStorage::Block::Node* a, BucketStorage::Block::Node* b)
{
	if (a != nullptr)
	{
		a->n_next = b;
		if (b != nullptr)
		{
			b->n_prev = a;
		}
	}
	else
	{
		if (b != nullptr)
		{
			b->n_prev = a;
		}
	}
}

template< typename T >
template< typename V >
typename BucketStorage< T >::iterator BucketStorage< T >::insertSimple(V&& value)
{
	if (freeBlocks == 0)
	{
		size_type b_id = b_lastBlock != nullptr ? b_lastBlock->b_id + 1 : 0;
		auto block = new Block(blockCapacity, std::forward< V >(value), b_id);
		block->n_countFilled++;
		iterator result = iterator(block, block->n_thisBegin);
		freeBlocks++;
		c_countBlocks++;
		b_firstEmptyBlock = b_lastEmptyBlock = block;
		if (b_lastBlock != nullptr)
		{
			b_lastBlock->b_nextBlock = block;
			block->b_prevBlock = b_lastBlock;
			b_lastBlock->n_thisEnd->n_next = block->n_thisBegin;
			block->n_thisBegin->n_prev = b_lastBlock->n_thisEnd;
		}
		else
		{
			b_firstBlock = block;
		}
		b_lastBlock = block;
		c_bucketSize++;
		new (*(&block->n_thisBegin->n_value)) value_type(std::forward< V >(value));
		return result;
	}
	iterator result = iterator(b_firstEmptyBlock, b_firstEmptyBlock->n_firstEmpty);
	new (*(&b_firstEmptyBlock->n_firstEmpty->n_value)) value_type(std::forward< V >(value));
	c_bucketSize++;
	b_firstEmptyBlock->n_countFilled++;
	if (blockCapacity == b_firstEmptyBlock->n_countFilled)
	{
		b_firstEmptyBlock->b_isFullBlock = true;
		freeBlocks--;
		b_firstEmptyBlock->n_firstEmpty = nullptr;
		if (b_firstEmptyBlock->b_prevFreeBlock != nullptr)
		{
			b_firstEmptyBlock->b_prevFreeBlock->b_nextFreeBlock = b_firstEmptyBlock->b_nextFreeBlock;
			if (b_firstEmptyBlock->b_nextFreeBlock != nullptr)
			{
				b_firstEmptyBlock->b_nextFreeBlock->b_prevFreeBlock = b_firstEmptyBlock->b_prevFreeBlock;
			}
			else
			{
				b_lastEmptyBlock = nullptr;
			}
		}
		else
		{
			if (b_firstEmptyBlock->b_nextFreeBlock != nullptr)
			{
				b_firstEmptyBlock->b_nextFreeBlock->b_prevFreeBlock = b_firstEmptyBlock->b_prevFreeBlock;
			}
			else
			{
				b_lastEmptyBlock = nullptr;
			}
		}
		b_firstEmptyBlock = b_firstEmptyBlock->b_nextFreeBlock;
	}
	else
	{
		b_firstEmptyBlock->n_firstEmpty = b_firstEmptyBlock->n_firstEmpty->n_next;
	}
	return result;
}

template< typename T >
typename BucketStorage< T >::iterator BucketStorage< T >::insert(const value_type& value)
{
	return insertSimple< const value_type& >(value);
}

template< typename T >
typename BucketStorage< T >::iterator BucketStorage< T >::insert(value_type&& value)
{
	return insertSimple< value_type&& >(std::move(value));
}

template< typename T >
void BucketStorage< T >::simpleErase(BucketStorage::const_iterator it)
{
	lincNodes(it.currentNode->n_prev, it.currentNode->n_next);
	lincNodes(it.currentBlock->n_thisEnd, it.currentNode);
	if (it.currentBlock->n_thisBegin == it.currentNode)
	{
		it.currentBlock->n_thisBegin = it.currentNode->n_next;
	}
	if (it.currentBlock->b_nextBlock != nullptr)
	{
		lincNodes(it.currentNode, it.currentBlock->b_nextBlock->n_thisBegin);
	}
	else
	{
		it.currentNode->n_next = nullptr;
	}
	it.currentBlock->n_thisEnd = it.currentNode;
	c_bucketSize--;
}

template< typename T >
typename BucketStorage< T >::iterator BucketStorage< T >::erase(const_iterator it)
{
	iterator res_it = it;
	res_it++;
	iterator result = iterator(res_it.currentBlock, res_it.currentNode);
	it.currentBlock->n_countFilled--;
	simpleErase(it);
	if (it.currentBlock->b_isFullBlock)
	{
		it.currentBlock->b_isFullBlock = false;
		if (freeBlocks)
		{
			b_lastEmptyBlock->b_nextFreeBlock = it.currentBlock;
			it.currentBlock->b_prevFreeBlock = b_lastEmptyBlock;
			it.currentBlock->b_nextFreeBlock = nullptr;
		}
		else
		{
			it.currentBlock->b_nextFreeBlock = nullptr;
			it.currentBlock->b_prevFreeBlock = nullptr;
			b_firstEmptyBlock = it.currentBlock;
		}
		b_lastEmptyBlock = it.currentBlock;
		freeBlocks++;
		it.currentBlock->n_firstEmpty = it.currentNode;
	}
	if (it.currentBlock->n_countFilled == 0)
	{
		freeBlocks--;
		c_countBlocks--;
		if (it.currentBlock->n_thisEnd->n_next != nullptr && it.currentBlock->b_nextBlock == nullptr)
		{
			delete it.currentBlock->n_thisEnd->n_next;
		}
		if (it.currentBlock->b_nextFreeBlock != nullptr)
		{
			it.currentBlock->b_nextFreeBlock->b_prevFreeBlock = it.currentBlock->b_prevFreeBlock;
			if (it.currentBlock->b_prevFreeBlock != nullptr)
			{
				it.currentBlock->b_prevFreeBlock->b_nextFreeBlock = it.currentBlock->b_nextFreeBlock;
			}
		}
		else
		{
			if (it.currentBlock->b_prevFreeBlock != nullptr)
			{
				it.currentBlock->b_prevFreeBlock->b_nextFreeBlock = it.currentBlock->b_nextFreeBlock;
			}
		}
		if (it.currentBlock->b_nextBlock != nullptr)
		{
			it.currentBlock->b_nextBlock->b_prevBlock = it.currentBlock->b_prevBlock;
			if (it.currentBlock->b_prevBlock != nullptr)
			{
				it.currentBlock->b_prevBlock->b_nextBlock = it.currentBlock->b_nextBlock;
			}
		}
		else
		{
			if (it.currentBlock->b_prevBlock != nullptr)
			{
				it.currentBlock->b_prevBlock->b_nextBlock = it.currentBlock->b_nextBlock;
			}
		}
		if (it.currentBlock == b_firstBlock)
		{
			b_firstBlock = it.currentBlock->b_nextBlock;
		}
		if (it.currentBlock == b_lastBlock)
		{
			b_lastBlock = it.currentBlock->b_prevBlock;
		}
		if (it.currentBlock == b_firstEmptyBlock)
		{
			b_firstEmptyBlock = it.currentBlock->b_nextFreeBlock;
		}
		if (it.currentBlock == b_lastEmptyBlock)
		{
			b_lastEmptyBlock = it.currentBlock->b_prevFreeBlock;
		}
		delete it.currentBlock;
	}
	return result;
}

template< typename T >
bool BucketStorage< T >::empty() const noexcept
{
	return c_countBlocks == 0;
}

template< typename T >
BucketStorage< T >::size_type BucketStorage< T >::size() const noexcept
{
	return c_bucketSize;
}

template< typename T >
BucketStorage< T >::size_type BucketStorage< T >::capacity() const noexcept
{
	return c_countBlocks * blockCapacity;
}

template< typename T >
void BucketStorage< T >::shrink_to_fit()
{
	BucketStorage< T > newStorage(blockCapacity);
	for (auto it = begin(); it != end(); ++it)
	{
		newStorage.insert(*it);
	}
	this->swap(newStorage);
}

template< typename T >
void BucketStorage< T >::clear()
{
	Block* block = b_firstBlock;
	if (b_lastBlock != nullptr && b_lastBlock->n_thisEnd != nullptr && b_lastBlock->n_thisEnd->n_next != nullptr)
	{
		delete b_lastBlock->n_thisEnd->n_next;
	}
	while (block != nullptr)
	{
		Block* nextBlock = block->b_nextBlock;
		delete block;
		block = nextBlock;
	}

	b_firstBlock = nullptr;
	b_lastBlock = nullptr;
	b_firstEmptyBlock = nullptr;
	b_lastEmptyBlock = nullptr;
	freeBlocks = 0;
	c_countBlocks = 0;
	c_bucketSize = 0;
}

template< typename T >
void BucketStorage< T >::swap(BucketStorage& other) noexcept
{
	using std::swap;
	swap(blockCapacity, other.blockCapacity);
	swap(c_countBlocks, other.c_countBlocks);
	swap(b_firstBlock, other.b_firstBlock);
	swap(b_lastBlock, other.b_lastBlock);
	swap(b_firstEmptyBlock, other.b_firstEmptyBlock);
	swap(b_lastEmptyBlock, other.b_lastEmptyBlock);
	swap(freeBlocks, other.freeBlocks);
	swap(c_bucketSize, other.c_bucketSize);
}

template< typename T >
typename BucketStorage< T >::iterator BucketStorage< T >::get_to_distance(iterator it, const difference_type distance)
{
	if (distance >= 0)
	{
		for (int i = 0; i < distance; i++)
		{
			it++;
		}
	}
	else
	{
		for (int i = 0; i < (-distance); i++)
		{
			it--;
		}
	}
	return it;
}

template< typename T >
typename BucketStorage< T >::iterator BucketStorage< T >::begin() noexcept
{
	return iterator(b_firstBlock, b_firstBlock == nullptr ? nullptr : b_firstBlock->n_thisBegin);
}
template< typename T >
typename BucketStorage< T >::const_iterator BucketStorage< T >::begin() const noexcept
{
	return const_iterator(b_firstBlock, b_firstBlock == nullptr ? nullptr : b_firstBlock->n_thisBegin);
}
template< typename T >
typename BucketStorage< T >::const_iterator BucketStorage< T >::cbegin() const noexcept
{
	return const_iterator(b_firstBlock, b_firstBlock == nullptr ? nullptr : b_firstBlock->n_thisBegin);
}
template< typename T >
typename BucketStorage< T >::iterator BucketStorage< T >::end() noexcept
{
	if (b_lastBlock != nullptr)
	{
		if (b_lastBlock->n_firstEmpty != nullptr)
		{
			return iterator(b_lastBlock, b_lastBlock->n_firstEmpty);
		}
		else
		{
			if (b_lastBlock->n_thisEnd->n_next != nullptr)
			{
				return iterator(b_lastBlock, b_lastBlock->n_thisEnd->n_next);
			}
			auto end_node = new Block::Node(b_lastBlock->n_thisEnd->n_id + 1);
			lincNodes(b_lastBlock->n_thisEnd, end_node);
			return iterator(b_lastBlock, end_node);
		}
	}
	else
	{
		return iterator(nullptr, nullptr);
	}
}
template< typename T >
typename BucketStorage< T >::const_iterator BucketStorage< T >::end() const noexcept
{
	if (b_lastBlock != nullptr)
	{
		if (b_lastBlock->n_firstEmpty != nullptr)
		{
			return const_iterator(b_lastBlock, b_lastBlock->n_firstEmpty);
		}
		else
		{
			if (b_lastBlock->n_thisEnd->n_next != nullptr)
			{
				return const_iterator(b_lastBlock, b_lastBlock->n_thisEnd->n_next);
			}
			auto end_node = new Block::Node(b_lastBlock->n_thisEnd->n_id + 1);
			lincNodes(b_lastBlock->n_thisEnd, end_node);
			return const_iterator(b_lastBlock, end_node);
		}
	}
	else
	{
		return const_iterator(nullptr, nullptr);
	}
}
template< typename T >
typename BucketStorage< T >::const_iterator BucketStorage< T >::cend() const noexcept
{
	if (b_lastBlock != nullptr)
	{
		if (b_lastBlock->n_firstEmpty != nullptr)
		{
			return const_iterator(b_lastBlock, b_lastBlock->n_firstEmpty);
		}
		else
		{
			if (b_lastBlock->n_thisEnd->n_next != nullptr)
			{
				return const_iterator(b_lastBlock, b_lastBlock->n_thisEnd->n_next);
			}
			auto end_node = new Block::Node(b_lastBlock->n_thisEnd->n_id + 1);
			b_lastBlock->n_thisEnd->n_next = end_node;
			end_node->n_prev = b_lastBlock->n_thisEnd;
			return const_iterator(b_lastBlock, end_node);
		}
	}
	else
	{
		return const_iterator(nullptr, nullptr);
	}
}

#endif	  // BUCKET_STORAGE_HPP