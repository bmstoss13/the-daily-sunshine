import type { Profile } from "./Profile";

const PostStatuses = {
    Draft: "draft",
    Published: "published",
    Archived: "archived"
} as const;

export type PostStatus = typeof PostStatuses[keyof typeof PostStatuses];

export interface Post {
    id: string;
    publisher_id: string;
    title: string;
    subtitle?: string;
    slug: string;
    content?: string;
    status: PostStatus;
    num_rays: number;
    num_comments: number;
    created_at: Date;
    updated_at: Date;
    deleted_at?: Date;

    publisher: Profile;
}



