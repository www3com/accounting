import {inject, observer} from "mobx-react";
import {Card, Col, Divider, Row, Space, Table} from "antd";

import Search from "@/pages/bill/component/Search";
import type {ColumnsType} from 'antd/es/table';

const main = ({store}: any) => {
    const columns: ColumnsType<any> = [
        {title: '日期', dataIndex: 'tradingTime', key: 'tradingTime', width: '120px', align: 'center'},
        {title: '业务类型', dataIndex: 'type', key: 'type', width: '80px', align: 'center'},
        {title: '分类/账户', dataIndex: 'account', key: 'account',},
        {title: '金额', dataIndex: 'amount', key: 'amount', align: 'right',},
        {title: '账户', dataIndex: 'cpAccount', key: 'cpAccount',},
        {title: '项目', dataIndex: 'project', key: 'project',},
        {title: '备注', dataIndex: 'remark', key: 'remark',}];

    const extra = (
        <Space split={<Divider type="vertical"/>}>
            <span>总支出 <span style={{color: '#14ba89'}}>-11,300.00</span></span>
            <span>总收入 <span style={{color: '#f1523a'}}>+29,100.00</span></span>
            <span>结余 17，800.00 <span style={{color: '#aaa'}}>（单位：元）</span></span>
        </Space>
    )

    return (
        <Card title='账目清单' size={'small'} bordered={false} extra={extra}>
            <Row gutter={[0, 8]}>
                <Col span={24}>
                    <Search/>
                </Col>
                <Col span={24}>
                    <Table size={"small"} rowKey={'id'} dataSource={store.transactions} bordered columns={columns}/>
                </Col>
            </Row>
        </Card>
    );
};

export default inject('store')(observer(main))
