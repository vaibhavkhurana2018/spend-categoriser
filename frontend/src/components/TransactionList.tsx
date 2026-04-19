import { useState, useEffect } from 'react';
import { useSearchParams } from 'react-router-dom';
import { useTransactions } from '../hooks/useTransactions';
import { db } from '../db';

export default function TransactionList() {
  const { transactions, categorySummary } = useTransactions();
  const [searchParams, setSearchParams] = useSearchParams();

  const [search, setSearch] = useState('');
  const [categoryFilter, setCategoryFilter] = useState(searchParams.get('category') ?? '');
  const [subCategoryFilter, setSubCategoryFilter] = useState(searchParams.get('subCategory') ?? '');
  const [page, setPage] = useState(0);
  const [sortField, setSortField] = useState<'date' | 'amount' | 'description'>('date');
  const [sortDir, setSortDir] = useState<'asc' | 'desc'>('desc');
  const perPage = 25;

  useEffect(() => {
    const cat = searchParams.get('category');
    const sub = searchParams.get('subCategory');
    if (cat) setCategoryFilter(cat);
    if (sub) setSubCategoryFilter(sub);
  }, [searchParams]);

  const subCategories = categoryFilter
    ? categorySummary.find((c) => c.name === categoryFilter)?.subCategories ?? []
    : [];

  const filtered = transactions
    .filter((t) => {
      if (search && !t.description.toLowerCase().includes(search.toLowerCase())) return false;
      if (categoryFilter && t.category !== categoryFilter) return false;
      if (subCategoryFilter && t.subCategory !== subCategoryFilter) return false;
      return true;
    })
    .sort((a, b) => {
      const dir = sortDir === 'asc' ? 1 : -1;
      if (sortField === 'amount') return (a.amount - b.amount) * dir;
      if (sortField === 'description') return a.description.localeCompare(b.description) * dir;
      return a.date.localeCompare(b.date) * dir;
    });

  const pageCount = Math.ceil(filtered.length / perPage);
  const paged = filtered.slice(page * perPage, (page + 1) * perPage);

  const toggleSort = (field: 'date' | 'amount' | 'description') => {
    if (sortField === field) setSortDir((d) => (d === 'asc' ? 'desc' : 'asc'));
    else { setSortField(field); setSortDir('desc'); }
  };

  const clearFilters = () => {
    setCategoryFilter('');
    setSubCategoryFilter('');
    setSearch('');
    setSearchParams({});
    setPage(0);
  };

  const hasFilters = categoryFilter || subCategoryFilter || search;

  const clearAll = async () => {
    if (confirm('Delete all transactions and bill history? This cannot be undone.')) {
      await db.transactions.clear();
      await db.bills.clear();
    }
  };

  if (transactions.length === 0) {
    return (
      <div className="text-center py-20">
        <div className="text-5xl mb-4">💳</div>
        <h2 className="text-xl font-bold text-gray-900">No transactions</h2>
        <p className="text-gray-500 mt-1">Upload credit card bills to see transactions</p>
      </div>
    );
  }

  const sortIcon = (field: string) => sortField === field ? (sortDir === 'asc' ? ' ↑' : ' ↓') : '';

  return (
    <div>
      <div className="flex items-center justify-between flex-wrap gap-2">
        <div>
          <h2 className="text-2xl font-bold text-gray-900">Transactions</h2>
          <p className="text-gray-500 mt-1">{filtered.length} of {transactions.length} transactions</p>
        </div>
        {hasFilters && (
          <button
            onClick={clearFilters}
            className="px-3 py-1.5 text-sm text-indigo-600 border border-indigo-200 rounded-lg hover:bg-indigo-50 transition-colors"
          >
            Clear Filters
          </button>
        )}
      </div>

      <div className="flex flex-wrap gap-3 mt-6">
        <input
          type="text"
          placeholder="Search transactions..."
          value={search}
          onChange={(e) => { setSearch(e.target.value); setPage(0); }}
          className="flex-1 min-w-[200px] px-4 py-2 border border-gray-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-transparent"
        />
        <select
          value={categoryFilter}
          onChange={(e) => { setCategoryFilter(e.target.value); setSubCategoryFilter(''); setPage(0); }}
          className="px-4 py-2 border border-gray-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-indigo-500"
        >
          <option value="">All Categories</option>
          {categorySummary.map((c) => (
            <option key={c.name} value={c.name}>{c.name} ({c.count})</option>
          ))}
        </select>
        {subCategories.length > 1 && (
          <select
            value={subCategoryFilter}
            onChange={(e) => { setSubCategoryFilter(e.target.value); setPage(0); }}
            className="px-4 py-2 border border-gray-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-indigo-500"
          >
            <option value="">All Subcategories</option>
            {subCategories.map((s) => (
              <option key={s.name} value={s.name}>{s.name} ({s.count})</option>
            ))}
          </select>
        )}
      </div>

      <div className="mt-4 bg-white rounded-xl border border-gray-200 overflow-x-auto">
        <table className="w-full text-sm">
          <thead className="bg-gray-50 border-b border-gray-200">
            <tr>
              <th onClick={() => toggleSort('date')} className="px-4 py-3 text-left font-medium text-gray-600 cursor-pointer hover:text-gray-900 whitespace-nowrap">Date{sortIcon('date')}</th>
              <th onClick={() => toggleSort('description')} className="px-4 py-3 text-left font-medium text-gray-600 cursor-pointer hover:text-gray-900">Description{sortIcon('description')}</th>
              <th onClick={() => toggleSort('amount')} className="px-4 py-3 text-right font-medium text-gray-600 cursor-pointer hover:text-gray-900 whitespace-nowrap">Amount{sortIcon('amount')}</th>
              <th className="px-4 py-3 text-left font-medium text-gray-600">Category</th>
              <th className="px-4 py-3 text-left font-medium text-gray-600">Sub-Category</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-gray-100">
            {paged.map((tx) => (
              <tr key={tx.id} className="hover:bg-gray-50 transition-colors">
                <td className="px-4 py-3 text-gray-600 whitespace-nowrap">{tx.date}</td>
                <td className="px-4 py-3 text-gray-900 font-medium">{tx.description}</td>
                <td className={`px-4 py-3 text-right font-medium whitespace-nowrap ${tx.amount < 0 ? 'text-green-600' : 'text-gray-900'}`}>
                  {tx.amount < 0 ? '−' : ''}{Math.abs(tx.amount).toLocaleString('en-IN', { style: 'currency', currency: 'INR', maximumFractionDigits: 2 })}
                </td>
                <td className="px-4 py-3"><span className="px-2 py-1 bg-indigo-50 text-indigo-700 rounded text-xs font-medium">{tx.category}</span></td>
                <td className="px-4 py-3 text-gray-500 text-xs">{tx.subCategory}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      {pageCount > 1 && (
        <div className="flex items-center justify-between mt-4">
          <button
            onClick={() => setPage((p) => Math.max(0, p - 1))}
            disabled={page === 0}
            className="px-4 py-2 text-sm border border-gray-300 rounded-lg disabled:opacity-50 hover:bg-gray-50 transition-colors"
          >Previous</button>
          <span className="text-sm text-gray-500">Page {page + 1} of {pageCount}</span>
          <button
            onClick={() => setPage((p) => Math.min(pageCount - 1, p + 1))}
            disabled={page === pageCount - 1}
            className="px-4 py-2 text-sm border border-gray-300 rounded-lg disabled:opacity-50 hover:bg-gray-50 transition-colors"
          >Next</button>
        </div>
      )}

      <div className="mt-8 pt-6 border-t border-gray-200">
        <button onClick={clearAll} className="px-4 py-2 text-sm text-red-600 border border-red-200 rounded-lg hover:bg-red-50 transition-colors">
          Clear All Data
        </button>
      </div>
    </div>
  );
}
