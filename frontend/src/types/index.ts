export interface Transaction {
  id?: number;
  date: string;
  description: string;
  amount: number;
  category: string;
  subCategory: string;
  hash: string;
  billFileName: string;
  addedAt: number;
}

export interface Bill {
  id?: number;
  fileName: string;
  uploadedAt: number;
  periodFrom: string;
  periodTo: string;
  totalAmount: number;
  transactionCount: number;
}

export interface ParseResponse {
  transactions: TransactionRaw[];
  billPeriod: { from: string; to: string };
  cardName: string;
  totalAmount: number;
  fileName: string;
}

export interface TransactionRaw {
  date: string;
  description: string;
  amount: number;
  category: string;
  subCategory: string;
  hash: string;
}

export interface CategorySummary {
  name: string;
  total: number;
  count: number;
  subCategories: { name: string; total: number; count: number }[];
}
