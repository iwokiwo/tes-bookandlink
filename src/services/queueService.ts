import axios from "axios";
import { baseUrl } from "../constants/form";

export interface QueuePayload {
  id?: string;
  email: string;
  url: string;
};

export interface QueueData {
  id: string;
  email: string;
  url: string;
  status: string;
  created_at: string;
  updated_at: string;
};

export const createQueueApi = async (
  payload: QueuePayload
) => {
  const response = await axios.post(
    `${baseUrl}/api/queue`,
    payload,
  );

  return response.data;
};

export const getQueueApi = async (): Promise<QueueData[]> => {
  try {
    const url = `${baseUrl}/api/queue`;
    const response = await axios.get<QueueData[]>(url);
    return response.data;
  } catch (error) {
    console.error("Failed to fetch queue data:", error);
    throw error;
  }
};

export const retryQueueApi = async (payload: QueuePayload) => {
  const { id, ...body } = payload;

  const response = await axios.post(
    `${baseUrl}/api/queue/${id}/retry`, // Sesuaikan dengan route API
    {...body,
      url: "http://www.google.com"
    },
  );

  return response.data;
};

export const deleteQueueApi = async () => {
  const response = await axios.delete(`${baseUrl}/api/queue`);
  return response.data;
};
