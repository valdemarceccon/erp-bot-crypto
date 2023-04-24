from typing import Dict
from typing import List

from pydantic import BaseModel


class TradeOrder(BaseModel):
    symbol: str
    orderId: str
    side: str
    qty: str
    orderPrice: str
    orderType: str
    execType: str
    closedSize: str
    cumEntryValue: str
    avgEntryPrice: str
    cumExitValue: str
    avgExitPrice: str
    closedPnl: str
    fillCount: str
    leverage: str
    createdTime: str
    updatedTime: str


class Result(BaseModel):
    list: List[TradeOrder]
    category: str
    nextPageCursor: str


class JsonResponse(BaseModel):
    retCode: int
    retMsg: str
    result: Result
    retExtInfo: Dict
    time: int


from src.dependencies.database import SessionLocal

from pybit import position
from pybit.unified_trading import HTTP


def get_closed_pnl(api_key, secret, days=1) -> JsonResponse:
    print(f"{api_key=} {secret=}")
    session = HTTP(
        testnet=False,
        api_key=api_key,
        api_secret=secret,
    )

    from datetime import datetime, timedelta

    now = datetime.now()
    yesterday = now - timedelta(days=days)

    start_time = yesterday.replace(hour=0, minute=0, second=0, microsecond=0)
    end_time = now.replace(hour=23, minute=59, second=59, microsecond=999999)

    start_timestamp = int(start_time.timestamp() * 1000)  # Convert to milliseconds
    end_timestamp = int(end_time.timestamp() * 1000)  # Convert to milliseconds

    r = session.get_closed_pnl(
        category="inverse",
        symbol="ETHUSD",
        startTime=start_timestamp,
        endTime=end_timestamp,
    )

    return JsonResponse.parse_obj(r)


def main():
    with SessionLocal() as db:
        from src.models.user import ClosedPnl, User, ApiKeyStatusEnum

        valdemar = (
            db.query(User).filter(User.email == "valdemar.ceccon@gmail.com").first()
        )
        if not valdemar:
            return
        for ak in valdemar.api_keys:
            if ak.status == ApiKeyStatusEnum.ACTIVE.value:
                for i in range(2, 16):
                    resp = get_closed_pnl(api_key=ak.api_key, secret=ak.secret, days=i)
                    all = [
                        ClosedPnl(user_id=valdemar.id, api_key_id=11, **r.dict())
                        for r in resp.result.list
                    ]
                    db.add_all(all)
        db.commit()


if __name__ == "__main__":
    main()
