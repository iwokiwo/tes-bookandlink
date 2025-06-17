import { createSlice, createAsyncThunk } from "@reduxjs/toolkit";
import { createQueueApi, deleteQueueApi, getQueueApi, retryQueueApi } from "../services/queueService";

export interface QueueFormData {
  email: string;
  url: string;
}

export interface QueueData {
  id: string;
  email: string;
  url: string;
  status: string;
  created_at: string;
  updated_at: string;
};

interface QueueState {
  loading: boolean;
  successMessage: string | null;
  errorMessage: string | null;
  queue: any,
  errorCode: number
}

const initialState: QueueState = {
  loading: false,
  successMessage: null,
  errorMessage: null,
  queue: [],
  errorCode: 200,
};

export interface CreateResponsesForm {
  answers: any[]
}

export const createQueueAsync = createAsyncThunk(
  "user/createQueue",
  async (formData: QueueFormData, { rejectWithValue }) => {
    try {
      const response = await createQueueApi(formData);
      return response;
    } catch (error: any) {
      if (error.response?.status === 422) {
        return rejectWithValue({ type: "fields", errors: error.response.data.errors });
      }
      if (error.response?.status === 401) {
        return rejectWithValue({ type: "auth", message: "You are not authorized." });
      }
      return rejectWithValue({ type: "general", message: "Something went wrong." });
    }
  }
);

export const retryQueueAsync = createAsyncThunk(
  "user/retryQueue",
  async ({
    formData
  }: { formData: QueueFormData }, { rejectWithValue }) => {
    try {

      const response = await retryQueueApi(formData);
      return response;
    } catch (error: any) {
      if (error.response?.status === 422) {
        return rejectWithValue({ type: "fields", errors: error.response.data.errors });
      }
      if (error.response?.status === 401) {
        return rejectWithValue({ type: "auth", message: "You are not authorized." });
      }
      return rejectWithValue({ type: "general", message: "Something went wrong." });
    }
  }
);

export const getQueueAsync = createAsyncThunk(
  "user/getQueue",
  async (_, { rejectWithValue }) => {
    try {
      const data = await getQueueApi();
      return data;
    } catch (error: any) {
      return rejectWithValue("Failed.");
    }
  }
);

export const deleteQueueAsync = createAsyncThunk(
  "user/deleteQueue",
  async (_, { rejectWithValue }) => {
    try {
      const data = await deleteQueueApi();
      return data;
    } catch (error: any) {
      return rejectWithValue("Failed.");
    }
  }
);


const queueSlice = createSlice({
  name: "queue",
  initialState,
  reducers: {
    clearQueueState: (state) => {
      state.loading = false;
      state.successMessage = null;
      state.errorMessage = null;
    },
  },
  extraReducers: (builder) => {
    builder
      .addCase(createQueueAsync.pending, (state) => {
        state.loading = true;
        state.successMessage = null;
        state.errorMessage = null;
        state.errorCode = 0
      })
      .addCase(createQueueAsync.fulfilled, (state, action) => {
        state.loading = false;
        state.successMessage = action.payload;
      })
      .addCase(createQueueAsync.rejected, (state, action: any) => {
        state.loading = false;
        state.errorMessage = action.payload?.message || "Failed";
      });
    builder
      .addCase(retryQueueAsync.pending, (state) => {
        state.loading = true;
        state.successMessage = null;
        state.errorMessage = null;
        state.errorCode = 0
      })
      .addCase(retryQueueAsync.fulfilled, (state, action) => {
        state.loading = false;
        state.successMessage = action.payload;
      })
      .addCase(retryQueueAsync.rejected, (state, action: any) => {
        state.loading = false;
        state.errorMessage = action.payload?.message || "Failed";
      });
    builder
      .addCase(getQueueAsync.pending, (state) => {
        state.loading = true;
        state.errorMessage = null;
        state.successMessage = null;
        state.errorCode = 0
      })
      .addCase(getQueueAsync.fulfilled, (state, action) => {
        state.loading = false;
        state.queue = action.payload;
      })
      .addCase(getQueueAsync.rejected, (state, action: any) => {
        state.loading = false;
        state.errorMessage = action.payload.message || "Failed";
        state.errorCode = action.payload?.code;
      });
      builder
      .addCase(deleteQueueAsync.pending, (state) => {
        state.loading = true;
        state.errorMessage = null;
        state.successMessage = null;
        state.errorCode = 0
      })
      .addCase(deleteQueueAsync.fulfilled, (state, action) => {
        state.loading = false;
        state.queue = action.payload;
      })
      .addCase(deleteQueueAsync.rejected, (state, action: any) => {
        state.loading = false;
        state.errorMessage = action.payload.message || "Failed";
        state.errorCode = action.payload?.code;
      });
  },
});

export const { clearQueueState } = queueSlice.actions;
export default queueSlice.reducer;
