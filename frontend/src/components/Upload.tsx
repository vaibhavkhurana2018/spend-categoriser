import { useState, useCallback } from 'react';
import { useDropzone } from 'react-dropzone';
import { db } from '../db';
import { deduplicateTransactions } from '../utils/dedup';
import type { ParseResponse } from '../types';

export default function Upload() {
  const [files, setFiles] = useState<File[]>([]);
  const [uploading, setUploading] = useState(false);
  const [results, setResults] = useState<{ fileName: string; added: number; duplicates: number }[]>([]);
  const [error, setError] = useState('');

  const onDrop = useCallback((accepted: File[]) => {
    setFiles((prev) => [...prev, ...accepted]);
    setError('');
    setResults([]);
  }, []);

  const { getRootProps, getInputProps, isDragActive } = useDropzone({
    onDrop,
    accept: { 'application/pdf': ['.pdf'] },
    multiple: true,
  });

  const removeFile = (index: number) => {
    setFiles((prev) => prev.filter((_, i) => i !== index));
  };

  const handleUpload = async () => {
    if (files.length === 0) return;
    setUploading(true);
    setError('');
    setResults([]);

    try {
      const formData = new FormData();
      files.forEach((f) => formData.append('files', f));

      const res = await fetch('/api/parse', { method: 'POST', body: formData });
      if (!res.ok) {
        const text = await res.text();
        throw new Error(text || 'Upload failed');
      }

      const parsed: ParseResponse[] = await res.json();
      const uploadResults: { fileName: string; added: number; duplicates: number }[] = [];

      for (const bill of parsed) {
        const { newTransactions, duplicateCount } = await deduplicateTransactions(
          bill.transactions,
          bill.fileName
        );

        if (newTransactions.length > 0) {
          await db.transactions.bulkAdd(newTransactions);
        }

        await db.bills.add({
          fileName: bill.fileName,
          uploadedAt: Date.now(),
          periodFrom: bill.billPeriod.from,
          periodTo: bill.billPeriod.to,
          totalAmount: bill.totalAmount,
          transactionCount: newTransactions.length,
        });

        uploadResults.push({
          fileName: bill.fileName,
          added: newTransactions.length,
          duplicates: duplicateCount,
        });
      }

      setResults(uploadResults);
      setFiles([]);
    } catch (err: any) {
      setError(err.message || 'Something went wrong');
    } finally {
      setUploading(false);
    }
  };

  return (
    <div>
      <h2 className="text-2xl font-bold text-gray-900">Upload Credit Card Bills</h2>
      <p className="text-gray-500 mt-1">Upload your PDF credit card statements to categorize spending</p>

      <div
        {...getRootProps()}
        className={`mt-6 border-2 border-dashed rounded-xl p-12 text-center cursor-pointer transition-colors ${
          isDragActive ? 'border-indigo-500 bg-indigo-50' : 'border-gray-300 hover:border-indigo-400 hover:bg-gray-50'
        }`}
      >
        <input {...getInputProps()} />
        <div className="text-5xl mb-4">📄</div>
        {isDragActive ? (
          <p className="text-indigo-600 font-medium">Drop your PDFs here...</p>
        ) : (
          <>
            <p className="text-gray-600 font-medium">Drag & drop PDF statements here</p>
            <p className="text-gray-400 text-sm mt-1">or click to browse files</p>
          </>
        )}
      </div>

      {files.length > 0 && (
        <div className="mt-6 bg-white rounded-xl border border-gray-200 divide-y">
          {files.map((f, i) => (
            <div key={i} className="flex items-center justify-between px-4 py-3">
              <div className="flex items-center gap-3">
                <span className="text-gray-400">📄</span>
                <span className="text-sm font-medium text-gray-700">{f.name}</span>
                <span className="text-xs text-gray-400">{(f.size / 1024).toFixed(0)} KB</span>
              </div>
              <button onClick={() => removeFile(i)} className="text-gray-400 hover:text-red-500 text-sm">Remove</button>
            </div>
          ))}
          <div className="p-4">
            <button
              onClick={handleUpload}
              disabled={uploading}
              className="w-full py-2.5 px-4 bg-indigo-600 text-white font-medium rounded-lg hover:bg-indigo-700 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
            >
              {uploading ? 'Processing...' : `Upload & Categorise ${files.length} file${files.length > 1 ? 's' : ''}`}
            </button>
          </div>
        </div>
      )}

      {error && (
        <div className="mt-4 p-4 bg-red-50 border border-red-200 rounded-xl">
          <p className="text-sm text-red-700">{error}</p>
        </div>
      )}

      {results.length > 0 && (
        <div className="mt-4 space-y-3">
          {results.map((r, i) => (
            <div key={i} className="p-4 bg-green-50 border border-green-200 rounded-xl">
              <p className="text-sm font-medium text-green-800">{r.fileName}</p>
              <p className="text-sm text-green-600 mt-1">
                {r.added} new transactions added
                {r.duplicates > 0 && `, ${r.duplicates} duplicates skipped`}
              </p>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
