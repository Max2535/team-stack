export type ApiError = {
  message: string;
  code: string;
};

export type ApiResponse<T> = {
  success: boolean;
  data?: T;
  error?: ApiError;
};

export type User = {
  id: string;
  email: string;
  name: string;
  role: string;
};
