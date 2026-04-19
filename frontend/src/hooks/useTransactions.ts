import { useLiveQuery } from 'dexie-react-hooks';
import { db } from '../db';
import type { CategorySummary } from '../types';

export function useTransactions() {
  const transactions = useLiveQuery(() => db.transactions.toArray()) ?? [];
  const bills = useLiveQuery(() => db.bills.orderBy('uploadedAt').reverse().toArray()) ?? [];

  const categorySummary: CategorySummary[] = (() => {
    const map = new Map<string, { total: number; count: number; subs: Map<string, { total: number; count: number }> }>();

    for (const tx of transactions) {
      if (tx.amount <= 0) continue;
      const cat = map.get(tx.category) ?? { total: 0, count: 0, subs: new Map() };
      cat.total += tx.amount;
      cat.count++;
      const sub = cat.subs.get(tx.subCategory) ?? { total: 0, count: 0 };
      sub.total += tx.amount;
      sub.count++;
      cat.subs.set(tx.subCategory, sub);
      map.set(tx.category, cat);
    }

    return Array.from(map.entries())
      .map(([name, data]) => ({
        name,
        total: Math.round(data.total * 100) / 100,
        count: data.count,
        subCategories: Array.from(data.subs.entries())
          .map(([subName, subData]) => ({
            name: subName,
            total: Math.round(subData.total * 100) / 100,
            count: subData.count,
          }))
          .sort((a, b) => b.total - a.total),
      }))
      .sort((a, b) => b.total - a.total);
  })();

  const totalSpend = transactions
    .filter((t) => t.amount > 0)
    .reduce((sum, t) => sum + t.amount, 0);

  const dateRange = (() => {
    if (transactions.length === 0) return { from: '', to: '' };
    const dates = transactions.map((t) => t.date).sort();
    return { from: dates[0], to: dates[dates.length - 1] };
  })();

  return { transactions, bills, categorySummary, totalSpend, dateRange };
}
