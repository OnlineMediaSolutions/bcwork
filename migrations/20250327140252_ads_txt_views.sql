-- +goose Up
-- +goose StatementBegin
create materialized view ads_txt_main_view as with main_table as (
    select t.*,
        p."name" as publisher_name,
        pd.mirror_publisher_id,
        pm."name" as mirror_publisher_name,
        p.account_manager_id,
        p.campaign_manager_id,
        u1.first_name || ' ' || u1.last_name as account_manager_full_name,
        u2.first_name || ' ' || u2.last_name as campaign_manager_full_name,
        u3.first_name || ' ' || u3.last_name as demand_manager_full_name
    from (
            select at2.id,
                at2.publisher_id,
                at2."domain",
                at2.domain_status,
                d.demand_partner_id,
                d.demand_partner_name,
                dpc.id as demand_partner_connection_id,
                dpc."media_type",
                d.demand_partner_name || ' - ' || d.demand_partner_name as demand_partner_name_extended,
                d.manager_id as demand_manager_id,
                at2.demand_status,
                at2.status,
                d.is_approval_needed,
                dpc.is_required_for_ads_txt as is_required,
                d.active as is_demand_partner_active,
                dpc.dp_domain || ', ' || dpc.publisher_account || ', ' || case
                    when dpc.is_direct then 'DIRECT'
                    else 'RESELLER'
                end || case
                    when dpc.certification_authority_id is not null then ', ' || dpc.certification_authority_id
                    else ''
                end as ads_txt_line,
                at2.last_scanned_at,
                at2.error_message
            from ads_txt at2
                join demand_partner_connection dpc on at2.demand_partner_connection_id = dpc.id
                join dpo d on d.demand_partner_id = dpc.demand_partner_id
            union
            select at2.id,
                at2.publisher_id,
                at2."domain",
                at2.domain_status,
                d.demand_partner_id,
                d.demand_partner_name as demand_partner_name,
                dpc2.id as demand_partner_connection_id,
                dpc2."media_type",
                d.demand_partner_name || ' - ' || dpc.dp_child_name as demand_partner_name_extended,
                d.manager_id as demand_manager_id,
                at2.demand_status,
                at2.status,
                d.is_approval_needed,
                dpc.is_required_for_ads_txt as is_required,
                d.active as is_demand_partner_active,
                dpc.dp_domain || ', ' || dpc.publisher_account || ', ' || case
                    when dpc.is_direct then 'DIRECT'
                    else 'RESELLER'
                end || case
                    when dpc.certification_authority_id is not null then ', ' || dpc.certification_authority_id
                    else ''
                end as ads_txt_line,
                at2.last_scanned_at,
                at2.error_message
            from ads_txt at2
                join demand_partner_child dpc on at2.demand_partner_child_id = dpc.id
                join demand_partner_connection dpc2 on dpc2.id = dpc.dp_connection_id
                join dpo d on d.demand_partner_id = dpc2.demand_partner_id
            union
            select at2.id,
                at2.publisher_id,
                at2."domain",
                at2.domain_status,
                '' as demand_partner_id,
                '' as demand_partner_name,
                0 as demand_partner_connection_id,
                null as media_type,
                so.seat_owner_name || ' - Direct' as demand_partner_name_extended,
                null as demand_manager_id,
                at2.demand_status,
                at2.status,
                false as is_approval_needed,
                true as is_required,
                true as is_demand_partner_active,
                so.seat_owner_domain || ', ' || replace(so.publisher_account, '%s', at2.publisher_id) || ', ' || 'DIRECT' || case
                    when so.certification_authority_id is not null then ', ' || so.certification_authority_id
                    else ''
                end as ads_txt_line,
                at2.last_scanned_at,
                at2.error_message
            from ads_txt at2
                join seat_owner so on at2.seat_owner_id = so.id
        ) as t
        join publisher p on p.publisher_id = t.publisher_id
        left join publisher_domain pd on pd.publisher_id = t.publisher_id
        and pd."domain" = t."domain"
        left join publisher pm on pm.publisher_id = pd.mirror_publisher_id
        left join "user" u1 on u1.id::varchar = p.account_manager_id
        left join "user" u2 on u2.id::varchar = p.campaign_manager_id
        left join "user" u3 on u3.id = t.demand_manager_id
)
select m1.id,
    m1.publisher_id,
    m1.publisher_name,
    m1.mirror_publisher_id,
    m1.mirror_publisher_name,
    m1.account_manager_full_name,
    m1.campaign_manager_full_name,
    m1.demand_manager_full_name,
    m1."domain",
    m1.demand_partner_id,
    m1.demand_partner_name,
    m1.demand_partner_name_extended,
    m1.demand_partner_connection_id,
    m1."media_type",
    m1.is_approval_needed,
    m1.is_required,
    m1.is_demand_partner_active,
    m1.last_scanned_at,
    m1.error_message,
    coalesce(m2.ads_txt_line, m1.ads_txt_line) as ads_txt_line,
    coalesce(m2.status, m1.status) as "status",
    coalesce(m2.domain_status, m1.domain_status) as domain_status,
    coalesce(m2.demand_status, m1.demand_status) as demand_status
from main_table m1
    left join main_table m2 on m1.mirror_publisher_id = m2.publisher_id
    and m1."domain" = m2."domain"
    and m1.demand_partner_name_extended = m2.demand_partner_name_extended
    and m1.demand_partner_connection_id = m2.demand_partner_connection_id;
--
create materialized view ads_txt_group_by_dp_view as with group_by_dp_table as (
select dense_rank() over (
        order by t.publisher_id,
            t."domain",
            t.demand_partner_name,
            t.demand_partner_connection_id
    ) as group_by_dp_id,
    t.*,
    p."name" as publisher_name,
    pd.mirror_publisher_id,
    pm."name" as mirror_publisher_name,
    p.account_manager_id,
    p.campaign_manager_id,
    u1.first_name || ' ' || u1.last_name as account_manager_full_name,
    u2.first_name || ' ' || u2.last_name as campaign_manager_full_name,
    u3.first_name || ' ' || u3.last_name as demand_manager_full_name,
    sum(
        case
            when t.status = 'added' then 1
            else 0
        end
    ) over (
        partition by t.publisher_id,
        t."domain",
        t.demand_partner_name,
        t.demand_partner_connection_id
    ) as added,
    count(t.status) over (
        partition by t.publisher_id,
        t."domain",
        t.demand_partner_name,
        t.demand_partner_connection_id
    ) as total,
    bool_and(
        case
            when t.status = 'added'
            AND t.is_required
            and t.demand_status = 'approved' then true
            when not t.is_required then true
            else false
        end
    ) over (
        partition by t.publisher_id,
        t."domain",
        t.demand_partner_name,
        t.demand_partner_connection_id
    ) as dp_enabled
from (
        select at2.id,
            at2.publisher_id,
            at2."domain",
            at2.domain_status,
            d.demand_partner_id,
            d.demand_partner_name,
            dpc.id as demand_partner_connection_id,
            dpc."media_type",
            d.demand_partner_name || ' - ' || d.demand_partner_name as demand_partner_name_extended,
            d.manager_id as demand_manager_id,
            at2.demand_status,
            at2.status,
            d.is_approval_needed,
            dpc.is_required_for_ads_txt as is_required,
            d.active as is_demand_partner_active,
            dpc.dp_domain || ', ' || dpc.publisher_account || ', ' || case
                when dpc.is_direct then 'DIRECT'
                else 'RESELLER'
            end || case
                when dpc.certification_authority_id is not null then ', ' || dpc.certification_authority_id
                else ''
            end as ads_txt_line,
            at2.last_scanned_at,
            at2.error_message
        from ads_txt at2
            join demand_partner_connection dpc on at2.demand_partner_connection_id = dpc.id
            join dpo d on d.demand_partner_id = dpc.demand_partner_id
        where dpc.is_required_for_ads_txt
        union
        select at2.id,
            at2.publisher_id,
            at2."domain",
            at2.domain_status,
            d.demand_partner_id,
            d.demand_partner_name as demand_partner_name,
            dpc2.id as demand_partner_connection_id,
            dpc2."media_type",
            d.demand_partner_name || ' - ' || dpc.dp_child_name as demand_partner_name_extended,
            d.manager_id as demand_manager_id,
            at2.demand_status,
            at2.status,
            d.is_approval_needed,
            dpc.is_required_for_ads_txt as is_required,
            d.active as is_demand_partner_active,
            dpc.dp_domain || ', ' || dpc.publisher_account || ', ' || case
                when dpc.is_direct then 'DIRECT'
                else 'RESELLER'
            end || case
                when dpc.certification_authority_id is not null then ', ' || dpc.certification_authority_id
                else ''
            end as ads_txt_line,
            at2.last_scanned_at,
            at2.error_message
        from ads_txt at2
            join demand_partner_child dpc on at2.demand_partner_child_id = dpc.id
            join demand_partner_connection dpc2 on dpc2.id = dpc.dp_connection_id
            join dpo d on d.demand_partner_id = dpc2.demand_partner_id
        union all
        select at2.id,
            at2.publisher_id,
            at2."domain",
            at2.domain_status,
            d.demand_partner_id,
            d.demand_partner_name as demand_partner_name,
            dpc.id as demand_partner_connection_id,
            dpc.media_type as media_type,
            so.seat_owner_name || ' - Direct' as demand_partner_name_extended,
            null as demand_manager_id,
            at2.demand_status,
            at2.status,
            d.is_approval_needed,
            true as is_required,
            d.active as is_demand_partner_active,
            so.seat_owner_domain || ', ' || replace(so.publisher_account, '%s', at2.publisher_id) || ', ' || 'DIRECT' || case
                when so.certification_authority_id is not null then ', ' || so.certification_authority_id
                else ''
            end as ads_txt_line,
            at2.last_scanned_at,
            at2.error_message
        from ads_txt at2
            join seat_owner so on at2.seat_owner_id = so.id
            join dpo d on d.seat_owner_id = so.id
            join demand_partner_connection dpc on d.demand_partner_id = dpc.demand_partner_id
    ) as t
    join publisher p on t.publisher_id = p.publisher_id
    left join publisher_domain pd on pd.publisher_id = t.publisher_id
    and pd."domain" = t."domain"
    left join publisher pm on pm.publisher_id = pd.mirror_publisher_id
    left join "user" u1 on u1.id::varchar = p.account_manager_id
    left join "user" u2 on u2.id::varchar = p.campaign_manager_id
    left join "user" u3 on u3.id = t.demand_manager_id
where t.is_demand_partner_active
)
select m1.group_by_dp_id as id,
    min(m1.publisher_id) as publisher_id,
    min(m1.publisher_name) as publisher_name,
    min(m1.mirror_publisher_id) as mirror_publisher_id,
    min(m1.mirror_publisher_name) as mirror_publisher_name,
    min(m1.account_manager_full_name) as account_manager_full_name,
    min(m1.campaign_manager_full_name) as campaign_manager_full_name,
    min(m1.demand_manager_full_name) as demand_manager_full_name,
    min(m1."domain") as "domain",
    min(coalesce(m2.domain_status, m1.domain_status)) as domain_status,
    min(m1.demand_partner_id) as demand_partner_id,
    min(m1.demand_partner_name) as demand_partner_name,
    min(m1.demand_partner_name_extended) as demand_partner_name_extended,
    min(m1.demand_partner_connection_id) as demand_partner_connection_id,
    min(m1."media_type") as "media_type",
    min(coalesce(m2.demand_status, m1.demand_status)) as demand_status,
    min(coalesce(m2.added, m1.added)) as added,
    min(coalesce(m2.total, m1.total)) as total,
    bool_and(coalesce(m2.dp_enabled, m1.dp_enabled)) as dp_enabled,
    min(m1.last_scanned_at) as last_scanned_at,
    json_agg(
        json_build_object(
            'id',
            m1.id,
            'demand_partner_name_extended',
            m1.demand_partner_name_extended,
            'demand_status',
            coalesce(m2.demand_status, m1.demand_status),
            'status',
            coalesce(m2.status, m1.status),
            'is_required',
            m1.is_required,
            'ads_txt_line',
            coalesce(m2.ads_txt_line, m1.ads_txt_line)
        )
    ) as grouped_lines_raw
from group_by_dp_table m1
    left join group_by_dp_table m2 on m1.mirror_publisher_id = m2.publisher_id
    and m1."domain" = m2."domain"
    and m1.demand_partner_name_extended = m2.demand_partner_name_extended
    and m1.demand_partner_connection_id = m2.demand_partner_connection_id
group by m1.group_by_dp_id;
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
drop materialized view ads_txt_main_view;
drop materialized view ads_txt_group_by_dp_view;
-- +goose StatementEnd