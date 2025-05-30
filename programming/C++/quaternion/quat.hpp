#include <cmath>
#include <iostream>

template< typename T >
struct matrix_t
{
	T data[16];
};

template< typename T >
struct vector3_t
{
	T x, y, z;
};

template< typename T >
class Quat
{
  public:
	Quat() : m_value{ 0, 0, 0, 0 } {}

	Quat(T a, T b, T c, T d) : m_value{ b, c, d, a } {}

	Quat(T angle, bool isRadians, vector3_t< T > axis)
	{
		T halfAngle = isRadians ? angle / 2 : angle * M_PI / 360;
		T invAxisLength = norm_vec(axis);
		T sinHalfAngle = std::sin(halfAngle);
		m_value[0] = axis.x * sinHalfAngle / invAxisLength;
		m_value[1] = axis.y * sinHalfAngle / invAxisLength;
		m_value[2] = axis.z * sinHalfAngle / invAxisLength;
		m_value[3] = std::cos(halfAngle);
	}

	Quat& operator+=(const Quat& other)
	{
		m_value[0] += other.m_value[0];
		m_value[1] += other.m_value[1];
		m_value[2] += other.m_value[2];
		m_value[3] += other.m_value[3];
		return *this;
	}
	Quat operator+(const Quat& other) const
	{
		Quat result(*this);
		result += other;
		return result;
	}
	Quat operator-(const Quat& other) const
	{
		Quat result(*this);
		result -= other;
		return result;
	}
	Quat operator*(const Quat& other) const
	{
		return Quat(
			(m_value[3] * other.m_value[3]) - (m_value[0] * other.m_value[0]) - (m_value[1] * other.m_value[1]) -
				(m_value[2] * other.m_value[2]),
			((m_value[3] * other.m_value[0]) + (other.m_value[3]) * m_value[0] + (m_value[1] * other.m_value[2]) -
			 (other.m_value[1] * m_value[2])),
			((m_value[3] * other.m_value[1]) + (other.m_value[3]) * m_value[1] + (m_value[2] * other.m_value[0]) -
			 (other.m_value[2] * m_value[0])),
			((m_value[3] * other.m_value[2]) + (other.m_value[3]) * m_value[2] + (m_value[0] * other.m_value[1]) -
			 (other.m_value[0] * m_value[1])));
	}
	Quat operator~() const { return Quat(m_value[3], -m_value[0], -m_value[1], -m_value[2]); }

	Quat& operator-=(const Quat& other)
	{
		m_value[0] -= other.m_value[0];
		m_value[1] -= other.m_value[1];
		m_value[2] -= other.m_value[2];
		m_value[3] -= other.m_value[3];
		return *this;
	}
	bool operator==(const Quat& other) const
	{
		return m_value[0] == other.m_value[0] && m_value[1] == other.m_value[1] && m_value[2] == other.m_value[2] &&
			   m_value[3] == other.m_value[3];
	}
	bool operator!=(const Quat& other) const { return !(*this == other); }
	Quat< T > operator*(const T& scalar) const
	{
		return Quat< T >(m_value[3] * scalar, m_value[0] * scalar, m_value[1] * scalar, m_value[2] * scalar);
	}
	Quat< T > operator*(const vector3_t< T >& vec) const { return Quat< T >(0, vec.x, vec.y, vec.z) * (*this); }
	explicit operator T() const
	{
		return std::sqrt(m_value[0] * m_value[0] + m_value[1] * m_value[1] + m_value[2] * m_value[2] + m_value[3] * m_value[3]);
	}

	matrix_t< T > rotation_matrix() const
	{
		matrix_t< T > result;
		T norm = static_cast< T >(*this);
		T norm4 = norm * norm;
		T qxx = m_value[0] * m_value[0] / norm4;
		T qyy = m_value[1] * m_value[1] / norm4;
		T qzz = m_value[2] * m_value[2] / norm4;
		T qxz = m_value[0] * m_value[2] / norm4;
		T qxy = m_value[0] * m_value[1] / norm4;
		T qyz = m_value[1] * m_value[2] / norm4;
		T qwx = m_value[3] * (-m_value[0]) / norm4;
		T qwy = m_value[3] * (-m_value[1]) / norm4;
		T qwz = m_value[3] * (-m_value[2]) / norm4;

		result.data[0] = 1 - 2 * (qyy + qzz);
		result.data[1] = 2 * (qxy - qwz);
		result.data[2] = 2 * (qxz + qwy);

		result.data[4] = 2 * (qxy + qwz);
		result.data[5] = 1 - 2 * (qxx + qzz);
		result.data[6] = 2 * (qyz - qwx);

		result.data[8] = 2 * (qxz - qwy);
		result.data[9] = 2 * (qyz + qwx);
		result.data[10] = 1 - 2 * (qxx + qyy);

		for (int i = 3; i < 12; i += 4)
		{
			result.data[i] = 0;
		}

		result.data[12] = result.data[13] = result.data[14] = 0;
		result.data[15] = 1;

		return result;
	}

	matrix_t< T > matrix() const
	{
		matrix_t< T > matrix1;
		matrix1.data[0] = matrix1.data[5] = matrix1.data[10] = matrix1.data[15] = m_value[3];
		matrix1.data[1] = matrix1.data[11] = -m_value[0];
		matrix1.data[4] = matrix1.data[14] = m_value[0];
		matrix1.data[2] = matrix1.data[13] = -m_value[1];
		matrix1.data[7] = matrix1.data[8] = m_value[1];
		matrix1.data[12] = matrix1.data[9] = m_value[2];
		matrix1.data[3] = matrix1.data[6] = -m_value[2];
		return matrix1;
	}

	T angle(bool degrees = true) const
	{
		T angle = 2 * std::acos(m_value[3]);
		return degrees ? angle : angle * 180 / M_PI;
	}

	vector3_t< T > apply(const vector3_t< T >& vec) const
	{
		T norm = static_cast< T >(*this);
		T norm4 = norm * norm * norm * norm;
		Quat q_norm = Quat();
		for (int i = 0; i < 4; ++i)
		{
			q_norm.m_value[i] = (m_value[i] / norm4);
		}
		vector3_t< T > vec_res;
		vec_res.x = -vec.x;
		vec_res.y = -vec.y;
		vec_res.z = -vec.z;
		Quat a = q_norm * vec_res;
		vector3_t< T > res;
		res.x = q_norm.m_value[0];
		res.y = q_norm.m_value[1];
		res.z = q_norm.m_value[2];
		return res;
	}

	const T* data() const { return m_value; }

  private:
	T m_value[4];

	T norm_vec(vector3_t< T >& axis) { return std::sqrt(axis.x * axis.x + axis.y * axis.y + axis.z * axis.z); }
};
