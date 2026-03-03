import { useCallback, useEffect, useRef, useState } from "react";
import { Link } from "react-router-dom";
import { toast } from "sonner";
import { profileSearch } from "../api/client";
import type { ProfileSummary } from "../api/types";
import styles from "./ProfileSearch.module.css";

export function ProfileSearch() {
  const [query, setQuery] = useState("");
  const [results, setResults] = useState<ProfileSummary[]>([]);
  const [loading, setLoading] = useState(false);
  const [searched, setSearched] = useState(false);
  const debounceRef = useRef<ReturnType<typeof setTimeout>>(undefined);

  const requestRef = useRef(0);

  const doSearch = useCallback(async (q: string) => {
    if (!q.trim()) {
      setResults([]);
      setSearched(false);
      return;
    }
    const seq = ++requestRef.current;
    setLoading(true);
    try {
      const res = await profileSearch(q);
      if (seq !== requestRef.current) return;
      setResults(res);
      setSearched(true);
    } catch (e) {
      if (seq !== requestRef.current) return;
      const msg = e instanceof Error ? e.message : "Search failed";
      toast.error(msg);
    } finally {
      if (seq === requestRef.current) setLoading(false);
    }
  }, []);

  useEffect(() => {
    if (debounceRef.current) clearTimeout(debounceRef.current);
    debounceRef.current = setTimeout(() => doSearch(query), 300);
    return () => {
      if (debounceRef.current) clearTimeout(debounceRef.current);
    };
  }, [query, doSearch]);

  return (
    <div className={styles.page}>
      <h2 className={styles.heading}>Profile Search</h2>
      <input
        className={styles.input}
        type="text"
        placeholder="Search by handle or name..."
        value={query}
        onChange={(e) => setQuery(e.target.value)}
        autoFocus
      />
      {loading && <p className={styles.dim}>Searching...</p>}
      {!loading && searched && results.length === 0 && <p className={styles.dim}>No profiles found.</p>}
      {results.length > 0 && (
        <ul className={styles.list}>
          {results.map((r) => (
            <li key={r.handle} className={styles.item}>
              <Link to={`/profile/${r.handle}`} className={styles.link}>
                <span className={styles.itemHandle}>@{r.handle}</span>
                {r.display_name && <span className={styles.itemName}>{r.display_name}</span>}
              </Link>
            </li>
          ))}
        </ul>
      )}
    </div>
  );
}
