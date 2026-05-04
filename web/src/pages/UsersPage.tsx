import { useEffect, useState } from "react";
import { useTranslation } from "react-i18next";
import { api } from "@/api";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table";
import { useToast } from "@/components/ui/toast";
import type { User } from "@/types";

function fmtTime(raw?: string) {
  if (!raw) return "-";
  const t = new Date(raw);
  if (Number.isNaN(t.getTime())) return raw;
  return t.toLocaleString();
}

export default function UsersPage() {
  const toast = useToast();
  const { t } = useTranslation();
  const [users, setUsers] = useState<User[]>([]);
  const [loading, setLoading] = useState(false);
  const [updatingId, setUpdatingId] = useState<string | null>(null);

  async function loadUsers() {
    setLoading(true);
    try {
      const res = await api.listUsers();
      setUsers(res.items || []);
    } catch (e) {
      const msg = e instanceof Error ? e.message : t("users.loadFailed");
      toast.error(msg);
    } finally {
      setLoading(false);
    }
  }

  useEffect(() => {
    void loadUsers();
  }, []);

  async function toggleStatus(user: User) {
    const next = user.status === "active" ? "disabled" : "active";
    setUpdatingId(user.id);
    try {
      await api.updateUserStatus(user.id, next);
      toast.success(t("users.updateSuccess"));
      await loadUsers();
    } catch (e) {
      const msg = e instanceof Error ? e.message : t("users.updateFailed");
      toast.error(msg);
    } finally {
      setUpdatingId(null);
    }
  }

  return (
    <div className="space-y-6">
      <Card>
        <CardHeader>
          <div className="flex items-center justify-between">
            <div>
              <CardTitle>{t("users.title")}</CardTitle>
              <CardDescription>{t("users.description")}</CardDescription>
            </div>
            <Button variant="outline" onClick={loadUsers} disabled={loading}>
              {t("common.refresh")}
            </Button>
          </div>
        </CardHeader>
        <CardContent>
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>{t("users.columns.email")}</TableHead>
                <TableHead>{t("users.columns.role")}</TableHead>
                <TableHead>{t("users.columns.status")}</TableHead>
                <TableHead>{t("users.columns.createdAt")}</TableHead>
                <TableHead className="text-right">{t("users.columns.actions")}</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {users.map((user) => (
                <TableRow key={user.id}>
                  <TableCell>{user.email}</TableCell>
                  <TableCell>
                    <Badge variant={user.role === "admin" ? "default" : "secondary"}>
                      {user.role === "admin" ? t("users.role.admin") : t("users.role.user")}
                    </Badge>
                  </TableCell>
                  <TableCell>
                    <Badge variant={user.status === "active" ? "secondary" : "destructive"}>
                      {user.status === "active" ? t("common.enabled") : t("common.disabled")}
                    </Badge>
                  </TableCell>
                  <TableCell>{fmtTime(user.createdAt)}</TableCell>
                  <TableCell className="text-right">
                    <Button
                      size="sm"
                      variant="outline"
                      disabled={updatingId === user.id || user.role === "admin"}
                      onClick={() => toggleStatus(user)}
                    >
                      {user.status === "active" ? t("common.disabled") : t("common.enabled")}
                    </Button>
                  </TableCell>
                </TableRow>
              ))}
              {users.length === 0 ? (
                <TableRow>
                  <TableCell colSpan={5} className="py-6 text-center text-sm text-[#666]">
                    {loading ? t("common.loading") : t("users.noUsers")}
                  </TableCell>
                </TableRow>
              ) : null}
            </TableBody>
          </Table>
        </CardContent>
      </Card>
    </div>
  );
}
