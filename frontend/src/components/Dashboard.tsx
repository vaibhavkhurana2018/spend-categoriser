import { useState, useMemo } from 'react';
import { Link } from 'react-router-dom';
import { useTransactions } from '../hooks/useTransactions';
import {
  PieChart, Pie, Cell, BarChart, Bar, XAxis, YAxis, Tooltip, ResponsiveContainer,
} from 'recharts';

const COLORS = ['#6366f1', '#8b5cf6', '#a78bfa', '#c084fc', '#e879f9', '#f472b6', '#fb7185', '#f87171', '#fb923c', '#fbbf24', '#a3e635', '#34d399'];

const MONTH_NAMES = ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec'];

function parseMonthKey(dateStr: string): string {
  const parts = dateStr.split('/');
  if (parts.length < 3) return '';
  const month = parseInt(parts[1], 10);
  const year = parts[2].length === 2 ? `20${parts[2]}` : parts[2];
  return `${year}-${String(month).padStart(2, '0')}`;
}

function parseYearKey(dateStr: string): string {
  const parts = dateStr.split('/');
  if (parts.length < 3) return '';
  return parts[2].length === 2 ? `20${parts[2]}` : parts[2];
}

function formatMonthLabel(key: string): string {
  const [year, month] = key.split('-');
  return `${MONTH_NAMES[parseInt(month, 10) - 1]} ${year}`;
}

function parseDateToSortable(d: string): string {
  const [day, month, year] = d.split('/');
  return `${year}/${month}/${day}`;
}

type PeriodType = 'all' | 'month' | 'year';

export default function Dashboard() {
  const { transactions } = useTransactions();

  const [selectedCategory, setSelectedCategory] = useState<string | null>(null);
  const [periodType, setPeriodType] = useState<PeriodType>('all');
  const [periodValue, setPeriodValue] = useState('');

  const { availableMonths, availableYears } = useMemo(() => {
    const monthSet = new Set<string>();
    const yearSet = new Set<string>();
    for (const tx of transactions) {
      const mk = parseMonthKey(tx.date);
      const yk = parseYearKey(tx.date);
      if (mk) monthSet.add(mk);
      if (yk) yearSet.add(yk);
    }
    return {
      availableMonths: Array.from(monthSet).sort(),
      availableYears: Array.from(yearSet).sort(),
    };
  }, [transactions]);

  const filteredTransactions = useMemo(() => {
    if (periodType === 'all') return transactions;
    return transactions.filter((tx) => {
      if (periodType === 'month') return parseMonthKey(tx.date) === periodValue;
      if (periodType === 'year') return parseYearKey(tx.date) === periodValue;
      return true;
    });
  }, [transactions, periodType, periodValue]);

  const categorySummary = useMemo(() => {
    const map = new Map<string, { total: number; count: number; subs: Map<string, { total: number; count: number }> }>();
    for (const tx of filteredTransactions) {
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
          .map(([subName, subData]) => ({ name: subName, total: Math.round(subData.total * 100) / 100, count: subData.count }))
          .sort((a, b) => b.total - a.total),
      }))
      .sort((a, b) => b.total - a.total);
  }, [filteredTransactions]);

  const totalSpend = useMemo(
    () => filteredTransactions.filter((t) => t.amount > 0).reduce((sum, t) => sum + t.amount, 0),
    [filteredTransactions],
  );

  const dateRange = useMemo(() => {
    if (filteredTransactions.length === 0) return { from: '', to: '' };
    const sorted = filteredTransactions.map((t) => t.date).sort((a, b) => parseDateToSortable(a).localeCompare(parseDateToSortable(b)));
    return { from: sorted[0], to: sorted[sorted.length - 1] };
  }, [filteredTransactions]);

  if (transactions.length === 0) {
    return (
      <div className="text-center py-20 animate-fade-in-up">
        <div className="text-6xl mb-4">📊</div>
        <h2 className="text-2xl font-bold text-gray-900">No transactions yet</h2>
        <p className="text-gray-500 mt-2">Upload your credit card statements to see spending insights</p>
        <Link to="/upload" className="inline-block mt-6 px-6 py-2.5 bg-indigo-600 text-white font-medium rounded-lg hover:bg-indigo-700 transition-colors">
          Upload Bills
        </Link>
      </div>
    );
  }

  const pieData = categorySummary.map((c) => ({ name: c.name, value: Math.round(c.total) }));
  const pieTotal = pieData.reduce((s, d) => s + d.value, 0);

  const topSubcategories = categorySummary
    .flatMap((c) => c.subCategories.map((s) => ({ name: s.name, total: s.total, category: c.name })))
    .sort((a, b) => b.total - a.total)
    .slice(0, 10);

  const monthlyMap = new Map<string, number>();
  for (const tx of filteredTransactions) {
    if (tx.amount <= 0) continue;
    const key = parseMonthKey(tx.date);
    if (key) monthlyMap.set(key, (monthlyMap.get(key) ?? 0) + tx.amount);
  }
  const monthlyTrend = Array.from(monthlyMap.entries())
    .sort(([a], [b]) => a.localeCompare(b))
    .map(([key, total]) => ({ month: formatMonthLabel(key), total: Math.round(total) }));

  const showTrend = monthlyTrend.length >= 2;

  const handleCategoryClick = (category: string) => {
    setSelectedCategory((prev) => (prev === category ? null : category));
  };

  const selectedTransactions = selectedCategory
    ? filteredTransactions
        .filter((tx) => tx.category === selectedCategory)
        .sort((a, b) => parseDateToSortable(b.date).localeCompare(parseDateToSortable(a.date)))
        .slice(0, 10)
    : [];

  const selectedCount = selectedCategory
    ? filteredTransactions.filter((tx) => tx.category === selectedCategory).length
    : 0;

  const handlePeriodChange = (type: PeriodType) => {
    setPeriodType(type);
    setSelectedCategory(null);
    if (type === 'all') {
      setPeriodValue('');
    } else if (type === 'month' && availableMonths.length > 0) {
      setPeriodValue(availableMonths[availableMonths.length - 1]);
    } else if (type === 'year' && availableYears.length > 0) {
      setPeriodValue(availableYears[availableYears.length - 1]);
    }
  };

  return (
    <div>
      <div className="flex items-center justify-between flex-wrap gap-2">
        <div>
          <h2 className="text-2xl font-bold text-gray-900">Dashboard</h2>
          <p className="text-gray-500 mt-1">Your spending overview</p>
        </div>

        {/* Period Filter */}
        <div className="flex items-center gap-2">
          <div className="inline-flex rounded-lg border border-gray-200 bg-white text-sm">
            {(['all', 'month', 'year'] as PeriodType[]).map((type) => (
              <button
                key={type}
                onClick={() => handlePeriodChange(type)}
                className={`px-3 py-1.5 font-medium transition-colors first:rounded-l-lg last:rounded-r-lg ${
                  periodType === type
                    ? 'bg-indigo-600 text-white'
                    : 'text-gray-600 hover:bg-gray-50'
                }`}
              >
                {type === 'all' ? 'All Time' : type === 'month' ? 'Month' : 'Year'}
              </button>
            ))}
          </div>
          {periodType === 'month' && (
            <select
              value={periodValue}
              onChange={(e) => { setPeriodValue(e.target.value); setSelectedCategory(null); }}
              className="px-3 py-1.5 border border-gray-200 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-indigo-500"
            >
              {availableMonths.map((m) => (
                <option key={m} value={m}>{formatMonthLabel(m)}</option>
              ))}
            </select>
          )}
          {periodType === 'year' && (
            <select
              value={periodValue}
              onChange={(e) => { setPeriodValue(e.target.value); setSelectedCategory(null); }}
              className="px-3 py-1.5 border border-gray-200 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-indigo-500"
            >
              {availableYears.map((y) => (
                <option key={y} value={y}>{y}</option>
              ))}
            </select>
          )}
        </div>
      </div>

      {filteredTransactions.length === 0 ? (
        <div className="text-center py-16 mt-6">
          <div className="text-4xl mb-3">📭</div>
          <p className="text-gray-500">No transactions in this period</p>
        </div>
      ) : (
        <>
          <div className="grid grid-cols-2 sm:grid-cols-4 gap-4 mt-6">
            {[
              { label: 'Total Spend', value: totalSpend.toLocaleString('en-IN', { style: 'currency', currency: 'INR', maximumFractionDigits: 0 }) },
              { label: 'Transactions', value: filteredTransactions.length },
              { label: 'Categories', value: categorySummary.length },
              { label: 'Period', value: `${dateRange.from} — ${dateRange.to}`, small: true },
            ].map((card, i) => (
              <div key={i} className="bg-white rounded-xl border border-gray-200 p-5 animate-fade-in-up" style={{ animationDelay: `${i * 80}ms`, animationFillMode: 'backwards' }}>
                <p className="text-sm text-gray-500">{card.label}</p>
                <p className={`${card.small ? 'text-lg' : 'text-2xl'} font-bold text-gray-900 mt-1`}>{card.value}</p>
              </div>
            ))}
          </div>

          <div className="grid grid-cols-1 lg:grid-cols-2 gap-6 mt-8">
            {/* Pie Chart */}
            <div className="bg-white rounded-xl border border-gray-200 p-6 animate-fade-in-up" style={{ animationDelay: '320ms', animationFillMode: 'backwards' }}>
              <h3 className="text-lg font-semibold text-gray-900 mb-4">Spend by Category</h3>
              <ResponsiveContainer width="100%" height={260}>
                <PieChart>
                  <Pie
                    data={pieData}
                    cx="50%"
                    cy="50%"
                    innerRadius={55}
                    outerRadius={100}
                    paddingAngle={2}
                    dataKey="value"
                    cursor="pointer"
                    onClick={(_, index) => handleCategoryClick(pieData[index].name)}
                  >
                    {pieData.map((d, i) => (
                      <Cell
                        key={i}
                        fill={COLORS[i % COLORS.length]}
                        opacity={selectedCategory && selectedCategory !== d.name ? 0.35 : 1}
                      />
                    ))}
                  </Pie>
                  <Tooltip formatter={(value: number) => value.toLocaleString('en-IN', { style: 'currency', currency: 'INR', maximumFractionDigits: 0 })} />
                </PieChart>
              </ResponsiveContainer>
              <div className="mt-3 max-h-40 overflow-y-auto grid grid-cols-2 gap-x-4 gap-y-1.5 text-xs">
                {pieData.map((d, i) => (
                  <button
                    key={d.name}
                    onClick={() => handleCategoryClick(d.name)}
                    className={`flex items-center gap-2 text-left rounded px-1 py-0.5 transition-colors ${
                      selectedCategory === d.name ? 'bg-indigo-50 ring-1 ring-indigo-200' : 'hover:bg-gray-50'
                    }`}
                  >
                    <span className="w-2.5 h-2.5 rounded-full shrink-0" style={{ backgroundColor: COLORS[i % COLORS.length] }} />
                    <span className="text-gray-700 truncate">{d.name}</span>
                    <span className="text-gray-400 ml-auto shrink-0">{pieTotal > 0 ? `${Math.round((d.value / pieTotal) * 100)}%` : ''}</span>
                  </button>
                ))}
              </div>
            </div>

            {/* Bar Chart */}
            <div className="bg-white rounded-xl border border-gray-200 p-6 animate-fade-in-up" style={{ animationDelay: '400ms', animationFillMode: 'backwards' }}>
              <h3 className="text-lg font-semibold text-gray-900 mb-4">Top Subcategories</h3>
              <ResponsiveContainer width="100%" height={300}>
                <BarChart
                  data={topSubcategories}
                  layout="vertical"
                  margin={{ left: 80 }}
                  onClick={(state) => {
                    if (state?.activePayload?.[0]?.payload) {
                      const item = state.activePayload[0].payload;
                      handleCategoryClick(item.category);
                    }
                  }}
                >
                  <XAxis type="number" tickFormatter={(v) => `₹${(v / 1000).toFixed(0)}k`} />
                  <YAxis type="category" dataKey="name" width={80} tick={{ fontSize: 12 }} />
                  <Tooltip formatter={(value: number) => value.toLocaleString('en-IN', { style: 'currency', currency: 'INR', maximumFractionDigits: 0 })} />
                  <Bar dataKey="total" fill="#6366f1" radius={[0, 4, 4, 0]} cursor="pointer" />
                </BarChart>
              </ResponsiveContainer>
            </div>
          </div>

          {/* Inline Transaction Panel */}
          {selectedCategory && (
            <div className="mt-6 bg-white rounded-xl border border-indigo-200 p-6 animate-fade-in-up">
              <div className="flex items-center justify-between mb-4">
                <div>
                  <h3 className="text-lg font-semibold text-gray-900">{selectedCategory}</h3>
                  <p className="text-sm text-gray-500">{selectedCount} transaction{selectedCount !== 1 ? 's' : ''}</p>
                </div>
                <button
                  onClick={() => setSelectedCategory(null)}
                  className="p-1.5 rounded-lg hover:bg-gray-100 transition-colors text-gray-400 hover:text-gray-600"
                  aria-label="Close"
                >
                  <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                  </svg>
                </button>
              </div>
              <div className="overflow-x-auto">
                <table className="w-full text-sm">
                  <thead className="bg-gray-50 border-b border-gray-200">
                    <tr>
                      <th className="px-4 py-2 text-left font-medium text-gray-600">Date</th>
                      <th className="px-4 py-2 text-left font-medium text-gray-600">Description</th>
                      <th className="px-4 py-2 text-right font-medium text-gray-600">Amount</th>
                      <th className="px-4 py-2 text-left font-medium text-gray-600">Sub-Category</th>
                    </tr>
                  </thead>
                  <tbody className="divide-y divide-gray-100">
                    {selectedTransactions.map((tx) => (
                      <tr key={tx.id} className="hover:bg-gray-50">
                        <td className="px-4 py-2 text-gray-600 whitespace-nowrap">{tx.date}</td>
                        <td className="px-4 py-2 text-gray-900 font-medium">{tx.description}</td>
                        <td className="px-4 py-2 text-right font-medium text-gray-900 whitespace-nowrap">
                          {Math.abs(tx.amount).toLocaleString('en-IN', { style: 'currency', currency: 'INR', maximumFractionDigits: 2 })}
                        </td>
                        <td className="px-4 py-2 text-gray-500 text-xs">{tx.subCategory}</td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
              {selectedCount > 10 && (
                <div className="mt-3 text-center">
                  <Link
                    to={`/transactions?category=${encodeURIComponent(selectedCategory)}`}
                    className="text-sm text-indigo-600 hover:text-indigo-700 font-medium"
                  >
                    View All {selectedCount} in Transactions →
                  </Link>
                </div>
              )}
            </div>
          )}

          {/* Monthly Trend */}
          {showTrend && (
            <div className="mt-8 bg-white rounded-xl border border-gray-200 p-6 animate-fade-in-up" style={{ animationDelay: '480ms', animationFillMode: 'backwards' }}>
              <h3 className="text-lg font-semibold text-gray-900 mb-4">Monthly Spending</h3>
              <ResponsiveContainer width="100%" height={280}>
                <BarChart data={monthlyTrend} margin={{ left: 20, right: 20, top: 10, bottom: 0 }}>
                  <XAxis dataKey="month" tick={{ fontSize: 12 }} />
                  <YAxis tickFormatter={(v) => `₹${(v / 1000).toFixed(0)}k`} tick={{ fontSize: 12 }} />
                  <Tooltip formatter={(value: number) => value.toLocaleString('en-IN', { style: 'currency', currency: 'INR', maximumFractionDigits: 0 })} />
                  <Bar dataKey="total" fill="#6366f1" radius={[4, 4, 0, 0]} />
                </BarChart>
              </ResponsiveContainer>
            </div>
          )}

          {/* Category Breakdown */}
          <div className="mt-8 bg-white rounded-xl border border-gray-200 p-6 animate-fade-in-up" style={{ animationDelay: showTrend ? '560ms' : '480ms', animationFillMode: 'backwards' }}>
            <h3 className="text-lg font-semibold text-gray-900 mb-4">Category Breakdown</h3>
            <div className="space-y-3">
              {categorySummary.map((cat, i) => (
                <div
                  key={cat.name}
                  className={`border rounded-lg p-4 cursor-pointer transition-colors ${
                    selectedCategory === cat.name
                      ? 'border-indigo-300 bg-indigo-50/50'
                      : 'border-gray-100 hover:border-indigo-200 hover:bg-indigo-50/30'
                  }`}
                  onClick={() => handleCategoryClick(cat.name)}
                >
                  <div className="flex items-center justify-between">
                    <div className="flex items-center gap-3">
                      <div className="w-3 h-3 rounded-full" style={{ backgroundColor: COLORS[i % COLORS.length] }} />
                      <span className="font-medium text-gray-900">{cat.name}</span>
                      <span className="text-xs text-gray-400">{cat.count} transactions</span>
                    </div>
                    <span className="font-semibold text-gray-900">{cat.total.toLocaleString('en-IN', { style: 'currency', currency: 'INR', maximumFractionDigits: 0 })}</span>
                  </div>
                  {cat.subCategories.length > 1 && (
                    <div className="mt-2 ml-6 space-y-1">
                      {cat.subCategories.map((sub) => (
                        <div key={sub.name} className="flex justify-between text-sm">
                          <span className="text-gray-500">{sub.name} ({sub.count})</span>
                          <span className="text-gray-700">{sub.total.toLocaleString('en-IN', { style: 'currency', currency: 'INR', maximumFractionDigits: 0 })}</span>
                        </div>
                      ))}
                    </div>
                  )}
                </div>
              ))}
            </div>
          </div>
        </>
      )}
    </div>
  );
}
