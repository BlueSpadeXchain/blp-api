import axios, { AxiosError } from 'axios';

// Types for API responses
interface User {
  id: string;
  userid: string;
  wallet_address: string;
  wallet_type: string;
  balance: number; // int64 in Go maps to number in TypeScript
  perp_balance: number;
  escrow_balance: number;
  stake_balance: number;
  frozen_balance: number;
  created_at: string;
}

interface ApiResponse<T> {
  data: T;
  status: number;
  message?: string;
}

// Define the base URL of your API
const BASE_URL = 'http://localhost:8080/api/user';

// Helper functions for API requests
const createUser = async (walletAddress: string, walletType: string): Promise<ApiResponse<User> | undefined> => {
  try {
    const response = await axios.post<ApiResponse<User>>(
      `${BASE_URL}`,
      {
        walletAddress,
        walletType,
      },
      {
        params: {
          query: 'get-user-by-address',
          "address": walletAddress,
          "type": walletType,
        },
      }
    );
    console.log('User created:', response.data);
    return response.data;
  } catch (error) {
    if (error instanceof AxiosError) {
      console.error('Error creating user:', error.response?.data || error.message);
    } else {
      console.error('Unexpected error:', error);
    }
    return undefined;
  }
};

const getUserById = async (userId: string): Promise<ApiResponse<User> | undefined> => {
  try {
    const response = await axios.get<ApiResponse<User>>(`${BASE_URL}`, {
      params: {
        query: 'get-user-by-id',
        'user-id': userId,
      },
    });
    console.log('User found by ID:', response.data);
    return response.data;
  } catch (error) {
    if (error instanceof AxiosError) {
      console.error('Error fetching user by ID:', error.response?.data || error.message);
    } else {
      console.error('Unexpected error:', error);
    }
    return undefined;
  }
};

const getUserByAddress = async (walletAddress: string, walletType: string): Promise<ApiResponse<User> | undefined> => {
  try {
    const response = await axios.get<ApiResponse<User>>(`${BASE_URL}`, {
      params: {
        query: 'get-user-by-address',
        "address": walletAddress,
        "type": walletType,
      },
    });
    console.log('User found by address:', response.data);
    return response.data;
  } catch (error) {
    if (error instanceof AxiosError) {
      console.error('Error fetching user by address:', error.response?.data || error.message);
    } else {
      console.error('Unexpected error:', error);
    }
    return undefined;
  }
};

// Test script
const runTests = async () => {
  // Test data
  const walletAddress = '0x1234567890abcdef1234567890abcdef12345678';
  const walletType = 'ecdsa';
  const userId = 'abcd1234';
  const nonExistentUserId = 'nonexistent1234';

  try {
    // 1. Create a user
    console.log('\nCreating a user...');
    const createdUser = await createUser(walletAddress, walletType);

    // 2. Find an existing user by ID
    console.log('\nFetching user by ID...');
    const userById = await getUserById(userId);

    // 3. Find an existing user by address
    console.log('\nFetching user by address...');
    const userByAddress = await getUserByAddress(walletAddress, walletType);

    // 4. Find a non-existent user by ID
    console.log('\nFetching non-existent user by ID...');
    const nonExistentUser = await getUserById(nonExistentUserId);

    return {
      createdUser,
      userById,
      userByAddress,
      nonExistentUser,
    };
  } catch (error) {
    console.error('Test execution failed:', error);
    return undefined;
  }
};

// Run the tests
runTests().then(() => console.log('Tests completed'));
