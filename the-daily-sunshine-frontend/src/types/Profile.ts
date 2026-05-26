const SubscriberTiers = {
    Member: "Member",
    Subscriber: "Subscriber",
} as const;

export type SubscriberTier = typeof SubscriberTiers[keyof typeof SubscriberTiers]

const AppRoles = {
    Admin: "Admin",
    User: "User",
} as const 

export type AppRole = typeof AppRoles[keyof typeof AppRoles]

export interface Profile {
    id: string;
    first_name: string;
    last_name: string;
    username: string;
    profile_image_url?: string;
    bio?: string;
    num_rays: number;
    subscriber_tier: SubscriberTier;
    app_role: AppRole;
    created_at: Date;
    updated_at: Date;
    deleted_at?: Date;
}